package service

import (
	"auth/common/logger"
	"auth/common/utils"
	"auth/config"
	"auth/models"
	"auth/pb"
	"auth/repository"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type UserServiceInterface interface {
	Register(user models.User) (int, error)
	LogIn(signInInfo models.SignInData) (*models.JWTTokenResponse, error)
	UpdateProfile(user *models.User) error
	ViewProfile(userID int) (*models.User, error)
	LogOut(accessToken string) error
	RequestSent(userID int, requestedID int) error
	RequestAccept(userID int, requestedID int) error
	ManageConnection(userID int, friendID int) error
	ViewFriends(userID int) ([]*models.User, error)
}

type UserService struct {
	repository repository.UserRepositoryInterface
	gRPCClient pb.AttachmentServiceClient
	redisRepo  repository.RedisRepositoryInterface
}

func NewUserService(gRPCClient pb.AttachmentServiceClient, repository repository.UserRepositoryInterface, redisRepo repository.RedisRepositoryInterface) UserServiceInterface {
	return &UserService{gRPCClient: gRPCClient, repository: repository, redisRepo: redisRepo}
}

func (service *UserService) Register(user models.User) (int, error) {
	err := user.Validate()
	if err != nil {
		return 0, err
	}

	respUser, err := service.repository.FindByEmail(user.Email)
	logger.LogError(respUser, " ", err)
	if err != nil {
		return 0, err
	} else if respUser != nil {
		return 0, errors.New("already registered user")
	}

	user.Password = utils.HashPassword(user.Password)
	return service.repository.Register(user)
}

func (service *UserService) LogIn(signInInfo models.SignInData) (*models.JWTTokenResponse, error) {
	logger.LogInfo(signInInfo.Email)
	err := signInInfo.Validate()
	if err != nil {
		return nil, err
	}

	respUser, err := service.repository.FindByEmail(signInInfo.Email)
	logger.LogError(respUser, " ", err)
	if err != nil {
		return nil, err
	} else if respUser == nil {
		return nil, errors.New("user not found")
	}

	err = utils.ComparePassword(respUser.Password, signInInfo.Password)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	accessToken, err := utils.CreateToken(config.Config.AccessTokenExpiresIn, respUser.ID, config.Config.AccessTokenPrivateKey)
	if err != nil {
		logger.LogError(err)
		return nil, errors.New(err.Error())
	}

	refreshToken, err := utils.CreateToken(config.Config.RefreshTokenExpiresIn, respUser.ID, config.Config.RefreshTokenPrivateKey)
	if err != nil {
		logger.LogError(err)
		return nil, errors.New(err.Error())
	}

	expirationTime := time.Now().Add(config.Config.RefreshTokenExpiresIn)

	resp := models.JWTTokenResponse{
		Token:     accessToken,
		Refresh:   refreshToken,
		ExpiredAt: expirationTime.Unix(),
		Email:     signInInfo.Email,
	}
	marshalJson, _ := json.Marshal(resp)
	_ = service.redisRepo.SetValue(context.Background(), "access_token:"+accessToken, string(marshalJson), config.Config.AccessTokenExpiresIn)

	return &resp, nil
}

func (service *UserService) UpdateProfile(user *models.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}
	_, err = service.ViewProfile(user.ID)
	if err != nil {
		return err
	}

	if user.Password != "" {
		user.Password = utils.HashPassword(user.Password)
	}

	requestAttachments := &pb.RequestAttachments{}
	err = service.repository.UpdateProfile(user)
	if err != nil {
		return err
	}

	tmpAttachment := pb.RequestAttachment{Name: user.ProfilePicName, Path: user.ProfilePicPath, SourceType: "user", SourceId: uint64(user.ID)}
	requestAttachments.Attachments = append(requestAttachments.Attachments, &tmpAttachment)

	_, err = service.gRPCClient.CreateMultiple(context.Background(), requestAttachments)
	if err != nil {
		logger.LogError(err)
		return err
	}

	if len(user.Websites) != 0 {
		err := service.repository.UpdateWebsites(user.Websites, user.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *UserService) ViewProfile(userID int) (*models.User, error) {

	field := strconv.Itoa(userID)
	data, err := service.redisRepo.GetSingleData(context.Background(), "user", field)
	if err != nil {
		logger.LogError(err)
	}

	if data != "" {
		var user models.User
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			logger.LogError(err)
			return nil, err
		}

		return &user, nil
	}

	user, err := service.repository.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("not found")
	}

	value := map[string]interface{}{}
	byteData, err := json.Marshal(user)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	value[field] = byteData
	err = service.redisRepo.Set(context.Background(), "user", value)
	if err != nil {
		logger.LogError(err)
	}

	params := &pb.FindAllRequestParams{SourceId: int64(userID), SourceType: "user"}
	gRPCAttachments, err := service.gRPCClient.FetchAll(context.Background(), params)
	if err != nil {
		return nil, err
	}

	attachments := []models.Attachment{}

	ln := len(gRPCAttachments.Attachments)
	if ln > 0 {
		for _, attattachment := range gRPCAttachments.Attachments {
			attachment := models.Attachment{Name: attattachment.Name, Path: attattachment.Path}
			attachments = append(attachments, attachment)
		}

		user.ProfilePicName = attachments[ln-1].Name
		user.ProfilePicPath = attachments[ln-1].Path
	}
	return user, nil
}

func (service *UserService) LogOut(accessToken string) error {
	return service.redisRepo.Delete(context.Background(), "access_token:"+accessToken, nil)
}

func (r *UserService) RequestSent(userID int, requestedID int) error {

	if userID == requestedID {
		return errors.New("cannot sent accept to myself")
	}

	_, err := r.ViewProfile(userID)
	if err != nil {
		return err
	}

	_, err = r.ViewProfile(requestedID)
	if err != nil {
		return err
	}

	err = r.repository.IsAlreadyRequestSent(userID, requestedID)
	if err != nil {
		return errors.New("already sent")
	}

	return r.repository.RequestSent(userID, requestedID)
}

func (r *UserService) RequestAccept(userID int, requestedID int) error {
	if userID == requestedID {
		return errors.New("cannot sent accept to myself")
	}

	_, err := r.ViewProfile(userID)
	if err != nil {
		return err
	}

	_, err = r.ViewProfile(requestedID)
	if err != nil {
		return err

	}

	err = r.repository.IsAlreadyRequestSent(userID, requestedID)
	if err == nil {
		return errors.New("no friend request")
	}

	err = r.repository.IsAlreadyRequestAccepter(userID, requestedID)
	if err != nil {
		return errors.New("no friend request")
	}
	return r.repository.RequestAccept(userID, requestedID)
}

func (r *UserService) ManageConnection(userID int, friendID int) error {

	err := r.repository.IsAlreadyRequestAccepter(userID, friendID)
	logger.LogError(err)
	if err == nil {
		return errors.New("you are not friend")
	}

	err = r.repository.ManageConnection(userID, friendID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserService) ViewFriends(userID int) ([]*models.User, error) {
	_, err := r.ViewProfile(userID)
	if err != nil {
		return nil, err
	}

	return r.repository.ViewFriends(userID)
}

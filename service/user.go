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
	"time"
)

type UserServiceInterface interface {
	Register(user models.User) (int, error)
	LogIn(signInInfo models.SignInData) (*models.JWTTokenResponse, error)
	UpdateProfile(user *models.User) error
	ViewProfile(userID int) (*models.User, error)
	LogOut(accessToken string) error
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
	respUser, err := service.repository.FindByEmail(signInInfo.Email)
	logger.LogError(respUser, " ", err)
	if err != nil {
		return nil, err
	} else if respUser == nil {
		return nil, errors.New("user not found")
	}

	err = utils.ComparePassword(respUser.Password, signInInfo.Password)
	if err != nil {
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
	// oldUser,err:=service.repository.FindByID(user.ID)
	// if err != nil {
	// 	return err
	// }
	requestAttachments := &pb.RequestAttachments{}
	service.repository.UpdateProfile(user)
	tmpAttachment := pb.RequestAttachment{Name: user.ProfilePicName, Path: user.ProfilePicPath, SourceType: "user", SourceId: uint64(user.ID)}
	requestAttachments.Attachments = append(requestAttachments.Attachments, &tmpAttachment)

	_, err := service.gRPCClient.CreateMultiple(context.Background(), requestAttachments)
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
	params := &pb.FindAllRequestParams{SourceId: int64(userID), SourceType: "user"}
	gRPCAttachments, err := service.gRPCClient.FetchAll(context.Background(), params)
	if err != nil {
		return nil, err
	}

	attachments := []models.Attachment{}

	for _, attattachment := range gRPCAttachments.Attachments {
		attachment := models.Attachment{Name: attattachment.Name, Path: attattachment.Path}
		attachments = append(attachments, attachment)
	}
	user, err := service.repository.FindByID(userID)
	user.ProfilePicName = attachments[0].Name
	user.ProfilePicPath = attachments[0].Path
	return user, nil
}

func (service *UserService) LogOut(accessToken string) error {
	return service.redisRepo.Delete(context.Background(), "access_token:"+accessToken, nil)
}

package service

import (
	"auth/common/logger"
	"auth/common/utils"
	"auth/models"
	"auth/repository"
	"errors"
)

type UserServiceInterface interface {
	Register(user models.User) error
	LogIn(signInInfo models.SignInData) (*models.User, error)
	UpdateProfile(user *models.User) error
	ViewProfile(userID int) (*models.User,error)
}

type UserService struct {
	repository repository.UserRepositoryInterface
}

func NewUserService(repository repository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{repository: repository}
}

func (service *UserService) Register(user models.User) error {
	respUser, err := service.repository.FindByEmail(user.Email)
	logger.LogError(respUser, " ", err)
	if err != nil {
		return err
	} else if respUser != nil {
		return errors.New("already registered user")
	}

	user.Password = utils.HashPassword(user.Password)
	return service.repository.Register(user)
}

func (service *UserService) LogIn(signInInfo models.SignInData) (*models.User, error) {
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

	return respUser, nil
}

func (service *UserService) UpdateProfile(user *models.User) error {
	// oldUser,err:=service.repository.FindByID(user.ID)
	// if err != nil {
	// 	return err
	// }

	if len(user.Websites) != 0 {
		err := service.repository.UpdateWebsites(user.Websites, user.ID)
		if err != nil {
			return err
		}
	}

	return service.repository.UpdateProfile(user)
}

func (service *UserService) ViewProfile(userID int) (*models.User,error) {
	return service.repository.FindByID(userID)
}

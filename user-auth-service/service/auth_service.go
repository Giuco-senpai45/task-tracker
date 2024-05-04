package service

import (
	"ua-service/http/requests"
	"ua-service/http/responses"
	"ua-service/jwt"
	"ua-service/models"
	"ua-service/utils"
	"ua-service/utils/errors"
	"ua-service/utils/log"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(dbAdapter *gorm.DB) *AuthService {
	models.New(dbAdapter)
	return &AuthService{
		db: dbAdapter,
	}
}

func (s *AuthService) Register(userRequest requests.CreateUserRequest) (responses.UserResponse, error) {
	user := models.User{
		FirstName: userRequest.Firstname,
		LastName:  userRequest.Lastname,
		Password:  userRequest.Password,
		Email:     userRequest.Email,
	}

	foundUser, err := models.FindUserByEmail(user.Email)
	if err == nil {
		log.Info("User already exists with email %v: %v", foundUser.Email, user)
		return responses.UserResponse{}, errors.ErrDuplicateEmail
	}

	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		return responses.UserResponse{}, err
	}
	user.Password = hashedPass

	dbErr := user.AddUser()
	errors.DBErrorCheck(dbErr)

	return responses.UserResponse{
		User: user,
	}, nil
}

func (s *AuthService) Login(userRequest requests.LoginRequest) (responses.LoginResponse, error) {
	user, err := models.FindUserByEmail(userRequest.Email)
	if err != nil {
		return responses.LoginResponse{}, errors.ErrNonExistingUser
	}
	log.Info("User found with email %v: %v", userRequest.Email, user)

	if !utils.ComparePassword(user.Password, userRequest.Password) || user.Email != userRequest.Email {
		return responses.LoginResponse{}, errors.ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(user.Email, int(user.ID))
	if err != nil {
		return responses.LoginResponse{}, err
	}

	return responses.LoginResponse{
		Token: token,
	}, nil
}

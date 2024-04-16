package services

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/golang-grpc-proxy/config"
	"github.com/golang-grpc-proxy/src/entities"
	"github.com/golang-grpc-proxy/src/storage"
	"github.com/golang-grpc-proxy/src/utils"
	"github.com/golang-grpc-proxy/src/utils/jwt"
)

var (
	ErrUserExist = errors.New("user already exist")
	ErrUserNotExist = errors.New("user does not exist")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidBody = errors.New("invalid body")
	ErrInternal = errors.New("internal server error")
	validate = validator.New(validator.WithRequiredStructEnabled())
)

type RegisterPdo struct {
	Nickname string 
	Email string 
	Phone string
	Password string
	FirstName string 
	Surname string 
	LastName string
	Study string
	Work string
}

type Api struct {
	storage *storage.Storage
	secret string
}

func NewApi(storage *storage.Storage, config *config.Config) Api {
	return Api{storage: storage, secret: config.Secret}
}

func (s *Api) Register(user entities.User) (int, error) {
	err := validate.Struct(user)

	if err != nil {
		return 0, ErrInvalidBody
	}
	
	isUnique, err := s.storage.IsUserUnique(user.Email, user.Phone)

	if !isUnique {
		return 0, ErrUserExist
	}

	if err != nil {
		return 0, ErrInternal
	}

	passwordHash, err := utils.CreatePasswordHash(user.Password)
	if err != nil {
		return 0, ErrInternal
	}
	user.Password = passwordHash

	user, err = s.storage.CreateUser(user)
	if err != nil {
		return 0, ErrInternal
	}

	return int(user.ID), nil
}

func (s *Api) AuthentificateUser(email, password string) error {
	user, err := s.storage.GetUserByEmail(email)

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrUserNotExist
		}
		return ErrInternal
	}

	if !utils.IsHashMatchPassword(user.Password, password) {
		return ErrInvalidCredentials
	}

	return nil
}

func (s *Api) GetToken(email, password string) (string, error) {
	err := s.AuthentificateUser(email, password)
	
	if err != nil {
		return "", err
	}

	token := jwt.New(email, s.secret)
	return token.ToString(), nil
}
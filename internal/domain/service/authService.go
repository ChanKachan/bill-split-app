package service

import (
	"errors"
	"github.com/ChanKachan/bill-split-app/internal/domain/entity/user"
	"github.com/ChanKachan/bill-split-app/internal/domain/repository"
	"github.com/ChanKachan/bill-split-app/internal/utils"
	"os"
	"strconv"
)

type AuthService interface {
	RegisterUser(userData user.User) (string, error)
	Auth(userData user.User) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Регистрация пользователя
func (as *authService) RegisterUser(userData user.User) (string, error) {
	var err error
	userId, err := as.userRepo.GetUserIdByLogin(userData.Login)
	if err != nil {
		return "", err
	}

	if userId != 0 {
		return "", errors.New("user already exists")
	}

	userData.Password, err = utils.HashPassword(userData.Password)
	if err != nil {
		return "", err
	}

	userData.Id, err = as.userRepo.CreateUser(userData)
	if err != nil {
		return "", err
	}

	jwtKey := os.Getenv("jwtSecretKey")

	token, err := utils.GenerateCentrifugeToken(strconv.Itoa(userData.Id), jwtKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Аутентификация пользователя
// Сервис проверяет входные параметры
// Выдает jwt токен, где есть параметры: userID, актуальность токена (2 недели)
// Обязательные данные для работы
// login string
// password string
func (as *authService) Auth(userData user.User) (string, error) {
	userOfDb, err := as.userRepo.GetUserByLogin(userData.Login)
	if err != nil {
		return "", err
	}

	if userOfDb == nil {
		return "", errors.New("user not found")
	}

	if userOfDb.Login != userData.Login || !utils.CheckPassword(userData.Password, userOfDb.Password) {
		return "", errors.New("invalid Login or Password")
	}

	jwtKey := os.Getenv("jwtSecretKey")

	token, err := utils.GenerateCentrifugeToken(strconv.Itoa(userOfDb.Id), jwtKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

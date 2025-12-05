package service

type AuthServiceInterface interface {
}

type AuthService struct{}

func NewAuthService() AuthServiceInterface {
	return &AuthService{}
}

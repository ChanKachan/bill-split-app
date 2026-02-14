package service

import (
	"bill-split/internal/domain/entity/auth"
	"bill-split/internal/repository"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
}

type authService struct {
	user repository.UserRepository
}

func NewAuthService(user repository.UserRepository) AuthService {
	return &authService{
		user: user,
	}
}

func (service *authService) Auth(c *gin.Context) {
	var authData auth.Auth
	if err := c.ShouldBind(authData); err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

}

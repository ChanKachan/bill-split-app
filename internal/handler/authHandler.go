package handler

import (
	"bill-split/internal/domain/entity/user"
	"bill-split/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	RegisterUser(c *gin.Context)
	Auth(c *gin.Context)
}

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{
		service: service,
	}
}

func (h *authHandler) RegisterUser(c *gin.Context) {
	data := user.User{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.service.RegisterUser(data)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
	})
	return
}

func (h *authHandler) Auth(c *gin.Context) {
	data := user.User{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.service.Auth(data)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
	})
	return
}

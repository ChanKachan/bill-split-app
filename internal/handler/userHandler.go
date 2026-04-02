package handler

import (
	"bill-split/internal/domain/entity/user"
	"bill-split/internal/domain/service"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type UserHandler interface {
	UpdateUserData(c *gin.Context)
	GetUserData(c *gin.Context)
}

type userHandler struct {
	userService service.UserHttpService
}

func NewUserHandler(userService service.UserHttpService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (uh *userHandler) UpdateUserData(c *gin.Context) {
	var userData user.User
	if err := c.ShouldBind(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userInfoAny, ok := c.Get("userID")
	if !ok {
		log.Println("UpdateUserData | user_info not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_info not found in context",
		})
		return
	}

	userId, ok := userInfoAny.(int)
	if !ok {
		log.Println("UpdateUserData | Failed to get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	userData.Id = userId

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := uh.userService.UpdateUser(ctx, userData)
	if err != nil {
		log.Printf("UpdateUserData | Failed user data %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func (uh *userHandler) GetUserData(c *gin.Context) {
	userInfoAny, ok := c.Get("userID")
	if !ok {
		log.Println("GetUserData | user_info not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_info not found in context",
		})
		return
	}

	userId, ok := userInfoAny.(int)
	if !ok {
		log.Println("GetUserData | Failed to get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	targetUserID := userId

	// Если передан параметр id в query, используем его (опционально)
	//if c.Query("id") != "" {
	//	id, err := strconv.Atoi(c.Query("id"))
	//	if err == nil && id > 0 {
	//		targetUserID = id
	//	}
	//}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userData, err := uh.userService.GetUser(ctx, targetUserID)
	if err != nil {
		log.Printf("GetUserData | Failed to get user data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user data",
		})
		return
	}

	if userData == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Не возвращаем пароль в ответе
	userData.Password = ""

	c.JSON(http.StatusOK, userData)
}

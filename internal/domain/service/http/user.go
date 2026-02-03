package http

import (
	"bill-split/internal/domain/entity/user"
	"bill-split/internal/repository"
	"bill-split/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

type UserHttpService interface {
	Create(c *gin.Context)
}
type userHttpService struct {
	userRepo repository.UserRepository
}

func NewUserHttpService(userRepo repository.UserRepository) UserHttpService {
	return &userHttpService{
		userRepo: userRepo,
	}
}

func (u *userHttpService) Create(c *gin.Context) {
	var user user.User
	if err := c.ShouldBind(&user); err != nil {
		log.Println("UserHttpService.Create err:", err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.Email == "" || user.Phone == "" || user.Login == "" {
		c.JSON(400, gin.H{
			"error": "user email or phone or login is empty",
		})
		return
	}

	passwordHash, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Password = passwordHash

	userId, err := u.userRepo.CreateUser(user)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"user_id": userId,
		"Code":    200,
	})

	return
}

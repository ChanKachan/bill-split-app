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
	GetUserById(c *gin.Context)
	UpdateUser(c *gin.Context)
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

func (u *userHttpService) UpdateUser(c *gin.Context) {
	var user user.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
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

	err = u.userRepo.UpdateUser(user)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
	})
	return
}

func (u *userHttpService) GetUserById(c *gin.Context) {
	var user user.User
	if err := c.ShouldBind(&user); err != nil {
		log.Println("UserHttpService.Create err:", err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.Id == 0 {
		c.JSON(400, gin.H{
			"error": "user id is empty",
		})
		return
	}

	userOld, err := u.userRepo.GetUserById(int(user.Id))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(200, userOld)
}

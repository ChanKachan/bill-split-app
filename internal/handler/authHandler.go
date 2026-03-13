package handler

import (
	"bill-split/internal/domain/entity/user"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) RegisterUser(c *gin.Context) {
	data := user.User{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.authorization.RegisterUser(data)
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

func (h *Handlers) Auth(c *gin.Context) {
	data := user.User{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.authorization.Auth(data)
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

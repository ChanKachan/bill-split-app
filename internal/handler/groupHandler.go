package handler

import (
	"bill-split/internal"
	"bill-split/internal/domain/entity/group"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (h *Handlers) CreateGroup(c *gin.Context) {
	groupInfo := group.Group{}
	if err := c.ShouldBind(&groupInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if groupInfo.Name == "" && groupInfo.DateStart.IsZero() && groupInfo.DateEnd.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group name and date cannot be both empty",
		})
		return
	}

	userInfoAny, ok := c.Get("user_info")
	if !ok {
		log.Println("Create Group | user_info not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_info not found in context",
		})
		return
	}

	userInfo, ok := userInfoAny.(*internal.UserInfo)
	if !ok {
		log.Println("Create Group | Failed to get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.groupService.CreateGroup(ctx, groupInfo, userInfo)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	return
}

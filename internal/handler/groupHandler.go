package handler

import (
	"bill-split/internal"
	"bill-split/internal/domain/entity/group"
	"bill-split/internal/domain/entity/groupMembers"
	"bill-split/internal/domain/service"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type GroupHandler interface {
	// POST
	CreateGroup(c *gin.Context)

	// GET
	GetUserGroups(c *gin.Context)
	GetUsersInGroup(c *gin.Context)
}

type groupHandlers struct {
	groupService service.GroupService
}

func NewGroupHandler(groupService service.GroupService) GroupHandler {
	return &groupHandlers{
		groupService: groupService,
	}
}

// Получить пользователей, которые находятся в группе
// Поиск введется по GroupId
func (h *groupHandlers) GetUsersInGroup(c *gin.Context) {
	groupMembersData := groupMembers.GroupMembers{}
	if err := c.ShouldBind(&groupMembersData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if groupMembersData.GroupId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group id can not be null",
		})
		return
	}

	userInfoAny, ok := c.Get("userID")
	if !ok {
		log.Println("GetUsersInGroup | user_info not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_info not found in context",
		})
		return
	}

	userId, ok := userInfoAny.(int)
	if !ok {
		log.Println("GetUsersInGroup | Failed to get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	userInfo := internal.UserInfo{UserId: userId}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	members, err := h.groupService.GetGroupMembersByGroupId(ctx, userInfo.UserId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": members,
	})
	return
}

// Создать группу
func (h *groupHandlers) CreateGroup(c *gin.Context) {
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

	userInfoAny, ok := c.Get("userID")
	if !ok {
		log.Println("Create Group | user_info not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_info not found in context",
		})
		return
	}

	userId, ok := userInfoAny.(int)
	if !ok {
		log.Println("Create Group | Failed to get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	userInfo := internal.UserInfo{UserId: userId}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.groupService.CreateGroup(ctx, groupInfo, &userInfo)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	return
}

func (gr *groupHandlers) GetUserGroups(c *gin.Context) {
	userInfoAny, ok := c.Get("userID")
	if !ok {
		log.Println("GetUserGroups | Failed get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get userInfoAny",
		})
		return
	}
	userId, ok := userInfoAny.(int)
	if !ok {
		log.Println("Create Group | Failed to get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	userInfo := internal.UserInfo{UserId: userId}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	groupsData, err := gr.groupService.GetGroupsById(ctx, userInfo.UserId)
	if err != nil {
		log.Printf("GetUserGroups | Failed to get groups by userId: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, groupsData)

	return
}

package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ChanKachan/bill-split-app/internal/domain/entity/group"
	"github.com/ChanKachan/bill-split-app/internal/domain/entity/groupMembers"
	"github.com/ChanKachan/bill-split-app/internal/domain/service"
	"github.com/ChanKachan/bill-split-app/internal/types"
	"github.com/gin-gonic/gin"
)

type GroupHandler interface {
	// POST
	CreateGroup(c *gin.Context)
	AddMember(c *gin.Context)
	LeaveGroup(c *gin.Context)
	RemoveUser(c *gin.Context) // Меняет статус del на 1
	AddUserToGroup(c *gin.Context)

	// GET
	GetUserGroups(c *gin.Context)
	GetUsersInGroup(c *gin.Context)
	GetGroupWithMembers(c *gin.Context)
}

type groupHandlers struct {
	groupService service.GroupService
}

func NewGroupHandler(groupService service.GroupService) GroupHandler {
	return &groupHandlers{
		groupService: groupService,
	}
}

func (h *groupHandlers) AddUserToGroup(c *gin.Context) {
	data := c.Param("link")
	if data == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("invalid data link"),
		})
		return
	}

	userInfoAny, ok := c.Get("userID")
	if !ok {
		log.Println("AddUserToGroup | user_info not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_info not found in context",
		})
		return
	}

	userId, ok := userInfoAny.(int)
	if !ok {
		log.Println("AddUserToGroup | Failed to get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	userInfo := types.UserInfo{UserId: userId}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	groupId, err := h.groupService.EnterUserToGroup(ctx, data, userInfo)
	if err != nil {
		log.Println("AddUserToGroup | Failed add user to group", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	groupData, members, err := h.groupService.GetGroupInfoWithMembers(ctx, groupId, userInfo)
	if err != nil {
		log.Println("AddUserToGroup | Failed get members with group info", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	dataRes := struct {
		Group   group.Group                 `json:"group"`
		Members []groupMembers.GroupMembers `json:"members"`
	}{
		Group:   groupData,
		Members: members,
	}

	c.JSON(http.StatusOK, dataRes)
}

func (h *groupHandlers) LeaveGroup(c *gin.Context) {
	groupMembersData := groupMembers.GroupMembers{}
	if err := c.ShouldBind(&groupMembersData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if groupMembersData.GroupId == 0 {
		log.Println("LeaveGroup | group id can not be empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group id can not be empty",
		})
		return
	}

	userInfoAny, ok := c.Get("userID")
	if !ok {
		log.Println("LeaveGroup | user_info not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_info not found in context",
		})
		return
	}

	userId, ok := userInfoAny.(int)
	if !ok {
		log.Println("LeaveGroup | Failed to get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.groupService.LeaveGroup(ctx, groupMembersData.GroupId, userId)
	if err != nil {
		log.Printf("LeaveGroup | Failed to leave group: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully left the group",
	})
}

func (h *groupHandlers) RemoveUser(c *gin.Context) {
	groupMembersData := groupMembers.GroupMembers{}
	if err := c.ShouldBind(&groupMembersData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if groupMembersData.GroupId == 0 || groupMembersData.UserId == 0 {
		log.Println("RemoveUser | group id/user id can not be empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group id/user id can not be empty",
		})
		return
	}

	userInfoAny, ok := c.Get("userID")
	if !ok {
		log.Println("RemoveUser | user_info not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_info not found in context",
		})
		return
	}

	requesterId, ok := userInfoAny.(int)
	if !ok {
		log.Println("RemoveUser | Failed to get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.groupService.RemoveUserFromGroup(ctx, groupMembersData.GroupId, groupMembersData.UserId, requesterId)
	if err != nil {
		log.Printf("RemoveUser | Failed to remove user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully removed from group",
	})
}

func (h *groupHandlers) AddMember(c *gin.Context) {
	groupMembersData := groupMembers.GroupMembers{}
	if err := c.ShouldBind(&groupMembersData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if groupMembersData.GroupId == 0 || groupMembersData.UserId == 0 || groupMembersData.Status == "" {
		log.Println("AddMember | group id/user id/status can not be empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group id/user id/status can not be empty",
		})
		return
	}

	userInfoAny, ok := c.Get("userID")
	if !ok {
		log.Println("AddMember | user_info not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_info not found in context",
		})
		return
	}

	userId, ok := userInfoAny.(int)
	if !ok {
		log.Println("AddMember | Failed to get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	userInfo := types.UserInfo{UserId: userId}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.groupService.AddUserToGroup(ctx, groupMembersData, userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	return
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

	//userInfoAny, ok := c.Get("userID")
	//if !ok {
	//	log.Println("GetUsersInGroup | user_info not found in context")
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": "user_info not found in context",
	//	})
	//	return
	//}

	//userId, ok := userInfoAny.(int)
	//if !ok {
	//	log.Println("GetUsersInGroup | Failed to get user info")
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": "Failed to get user info",
	//	})
	//	return
	//}
	//
	//userInfo := types.UserInfo{UserId: userId}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	members, err := h.groupService.GetGroupMembersByGroupId(ctx, 1) // todo: замоканные данные
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

	userInfo := types.UserInfo{UserId: userId}

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

	userInfo := types.UserInfo{UserId: userId}

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

func (gr *groupHandlers) GetGroupWithMembers(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id can not be empty",
		})
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	userInfoAny, ok := c.Get("userID")
	if !ok {
		log.Println("GetGroupWithMembers | Failed get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get userInfoAny",
		})
		return
	}
	userId, ok := userInfoAny.(int)
	if !ok {
		log.Println("GetGroupWithMembers | Failed to get user info")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	userInfo := types.UserInfo{UserId: userId}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	groupsData, members, err := gr.groupService.GetGroupInfoWithMembers(ctx, idInt, userInfo)
	if err != nil {
		log.Printf("GetGroupWithMembers | Failed to get groups by userId: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	data := struct {
		Group   group.Group                 `json:"group"`
		Members []groupMembers.GroupMembers `json:"members"`
	}{
		Group:   groupsData,
		Members: members,
	}

	c.JSON(http.StatusOK, data)

	return
}

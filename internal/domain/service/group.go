package service

import (
	"bill-split/internal"
	groupStruct "bill-split/internal/domain/entity/group"
	"bill-split/internal/domain/entity/groupMembers"
	"bill-split/internal/repository"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GroupService interface {
	CreateGroup(c *gin.Context)
	AddUserToGroup(c *gin.Context)
}
type group struct {
	groupRepo repository.GroupRepository
}

func NewGroupService(groupRepo repository.GroupRepository) GroupService {
	return &group{
		groupRepo: groupRepo,
	}
}

func (gr *group) CreateGroup(c *gin.Context) {
	groupInfo := groupStruct.Group{}

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

	err := gr.groupRepo.TransactionBegin(ctx)
	if err != nil {
		log.Println("Create Group | Failed begin transaction")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer func() {
		if err != nil {
			err := gr.groupRepo.TransactionRollback(ctx)
			if err != nil {
				log.Println("Create Group | Failed rollback transaction")
				log.Println(err.Error())
			}
		}

		err = gr.groupRepo.Commit(ctx)
		if err != nil {
			log.Println("Create Group | Failed commit group")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}()

	groupInfo.CreateAt = time.Now()
	groupInfo.Id, err = gr.groupRepo.CreateGroup(groupInfo)
	if err != nil {
		log.Println("Create Group | Failed create group")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	groupData := groupMembers.GroupMembers{
		UserId:  userInfo.UserId,
		GroupId: groupInfo.Id,
		Status:  string(groupMembers.Admin),
		Del:     0,
	}

	err = gr.groupRepo.AddUserToGroup(groupData)
	if err != nil {
		log.Println("Create Group | Failed add user to group")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	return
}

func (gr *group) AddUserToGroup(c *gin.Context) {
	var groupMembers groupMembers.GroupMembers

	if err := c.ShouldBind(&groupMembers); err != nil {
		log.Println("AddUserToGroup | Failed bind group members struct")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if err := gr.groupRepo.AddUserToGroup(groupMembers); err != nil {
		log.Println("AddUserToGroup | Failed add user to group")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	return
}

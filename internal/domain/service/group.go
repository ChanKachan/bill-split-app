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
	CreateGroup(ctx context.Context, groupInfo groupStruct.Group, userInfo *internal.UserInfo) error
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

func (gr *group) CreateGroup(ctx context.Context, groupInfo groupStruct.Group, userInfo *internal.UserInfo) error {

	err := gr.groupRepo.TransactionBegin(ctx)
	if err != nil {
		log.Println("Create Group | Failed begin transaction")
		return err
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
			log.Println(err.Error())
		}
	}()

	groupInfo.CreateAt = time.Now()
	groupInfo.Id, err = gr.groupRepo.CreateGroup(groupInfo)
	if err != nil {
		log.Println("Create Group | Failed create group")
		return err
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
		return err
	}
	return nil
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

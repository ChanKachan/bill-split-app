package service

import (
	"bill-split/internal"
	groupStruct "bill-split/internal/domain/entity/group"
	"bill-split/internal/domain/entity/groupMembers"
	"bill-split/internal/repository"
	"context"
	"errors"
	"log"
	"time"
)

type GroupService interface {
	// POST
	CreateGroup(ctx context.Context, groupInfo groupStruct.Group, userInfo *internal.UserInfo) error

	// GET
	GetGroupsById(ctx context.Context, userId int) ([]groupStruct.Group, error)
	GetGroupMembersByGroupId(ctx context.Context, groupId int) ([]groupMembers.GroupMembers, error)
}
type group struct {
	groupRepo repository.GroupRepository
}

func NewGroupService(groupRepo repository.GroupRepository) GroupService {
	return &group{
		groupRepo: groupRepo,
	}
}

// Создать событие
// Обязательные параметры
/*
	"name":"test", - название события
    "date_start": "2022-01-01T00:00:00Z", - дата начало события
    "date_end": "2022-01-01T00:00:00Z", - дата окончание события
    "amount": 12000 - потраченная сумма во время события
*/
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

//// Добавить пользователя в группу
//func (gr *group) AddUserToGroup(c *gin.Context) {
//	var groupMembers groupMembers.GroupMembers
//
//	if err := c.ShouldBind(&groupMembers); err != nil {
//		log.Println("AddUserToGroup | Failed bind group members struct")
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//
//	if err := gr.groupRepo.AddUserToGroup(groupMembers); err != nil {
//		log.Println("AddUserToGroup | Failed add user to group")
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": err.Error(),
//		})
//	}
//
//	return
//}

func (gr *group) GetGroupsById(ctx context.Context, userId int) ([]groupStruct.Group, error) {
	if userId == 0 {
		return nil, errors.New("userId can not be 0")
	}

	groupsData, err := gr.groupRepo.GetGroupsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return groupsData, nil
}

func (gr *group) GetGroupMembersByGroupId(ctx context.Context, groupId int) ([]groupMembers.GroupMembers, error) {
	if groupId == 0 {
		return nil, errors.New("group id can not be 0")
	}

	groupMembersData, err := gr.groupRepo.GetMembersByGroupId(ctx, groupId)
	if err != nil {
		return nil, err
	}

	return groupMembersData, nil
}

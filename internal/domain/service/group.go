package service

import (
	"context"
	"errors"
	"github.com/ChanKachan/bill-split-app/internal"
	groupStruct "github.com/ChanKachan/bill-split-app/internal/domain/entity/group"
	"github.com/ChanKachan/bill-split-app/internal/domain/entity/groupMembers"
	"github.com/ChanKachan/bill-split-app/internal/domain/repository"
	"log"
	"time"
)

type GroupService interface {
	// POST
	CreateGroup(ctx context.Context, groupInfo groupStruct.Group, userInfo *internal.UserInfo) error
	AddUserToGroup(ctx context.Context, member groupMembers.GroupMembers, userInfo internal.UserInfo) error
	LeaveGroup(ctx context.Context, groupId, userId int) error
	RemoveUserFromGroup(ctx context.Context, groupId, userId, requesterId int) error
	EnterUserToGroup(ctx context.Context, link string, userInfo internal.UserInfo) (int, error)

	// GET
	GetGroupsById(ctx context.Context, userId int) ([]groupStruct.Group, error)
	GetGroupMembersByGroupId(ctx context.Context, groupId int) ([]groupMembers.GroupMembers, error)
	CheckUserInGroup(ctx context.Context, groupId, userId int) (bool, error)
	GetGroupInfoWithMembers(ctx context.Context, groupId int, userInfo internal.UserInfo) (groupStruct.Group, []groupMembers.GroupMembers, error)
}
type group struct {
	groupRepo repository.GroupRepository
}

func NewGroupService(groupRepo repository.GroupRepository) GroupService {
	return &group{
		groupRepo: groupRepo,
	}
}

func (g *group) EnterUserToGroup(ctx context.Context, link string, userInfo internal.UserInfo) (int, error) {
	if link == "" {
		return 0, errors.New("link is empty")
	}

	groupid, err := g.groupRepo.GetGroupIdByLink(link)
	if err != nil {
		return 0, err
	}

	exist, err := g.groupRepo.CheckUserInGroup(ctx, groupid, userInfo.UserId)
	if err != nil {
		return 0, err
	}

	if exist {
		return 0, nil
	}

	err = g.AddUserToGroup(ctx, groupMembers.GroupMembers{GroupId: groupid, UserId: userInfo.UserId}, userInfo)
	if err != nil {
		return 0, err
	}
	return groupid, nil
}

func (g *group) GetGroupInfoWithMembers(ctx context.Context, groupId int, userInfo internal.UserInfo) (groupStruct.Group, []groupMembers.GroupMembers, error) {
	if groupId == 0 {
		return groupStruct.Group{}, nil, errors.New("invalid group id")
	}

	groupData, err := g.groupRepo.GetGroup(groupId)
	if err != nil {
		return groupStruct.Group{}, nil, err
	}

	members, err := g.groupRepo.GetMembersByGroupId(ctx, groupId)
	if err != nil {
		return groupStruct.Group{}, nil, err
	}

	return groupData, members, nil
}

func (gr *group) LeaveGroup(ctx context.Context, groupId, userId int) error {
	if groupId == 0 || userId == 0 {
		return errors.New("group id and user id cannot be 0")
	}

	// Проверяем, состоит ли пользователь в группе
	exists, err := gr.CheckUserInGroup(ctx, groupId, userId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user is not a member of this group")
	}

	// Проверяем, не является ли пользователь единственным админом
	role, err := gr.groupRepo.GetUserRoleInGroup(ctx, groupId, userId)
	if err != nil {
		return err
	}

	if role == string(groupMembers.Admin) {
		// Проверяем, есть ли другие админы в группе
		members, err := gr.groupRepo.GetMembersByGroupId(ctx, groupId)
		if err != nil {
			return err
		}

		adminCount := 0
		for _, member := range members {
			if member.Status == string(groupMembers.Admin) {
				adminCount++
			}
		}

		// Если это единственный админ, не даём выйти
		if adminCount <= 1 {
			return errors.New("cannot leave group: you are the only admin. Transfer admin rights to another user or delete the group")
		}
	}

	// Удаляем пользователя из группы (мягкое удаление)
	err = gr.groupRepo.RemoveUserFromGroup(ctx, groupId, userId)
	if err != nil {
		log.Printf("LeaveGroup | Failed to remove user from group: %v", err)
		return err
	}

	return nil
}

func (gr *group) RemoveUserFromGroup(ctx context.Context, groupId, userId, requesterId int) error {
	if groupId == 0 || userId == 0 || requesterId == 0 {
		return errors.New("group id, user id and requester id cannot be 0")
	}

	// Проверяем права запрашивающего
	requesterRole, err := gr.groupRepo.GetUserRoleInGroup(ctx, groupId, requesterId)
	if err != nil {
		return err
	}
	if requesterRole != string(groupMembers.Admin) {
		return errors.New("only admins can remove users from group")
	}

	// Проверяем, что удаляемый пользователь состоит в группе
	exists, err := gr.CheckUserInGroup(ctx, groupId, userId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user is not a member of this group")
	}

	// Не даём админу удалить самого себя через этот метод (использовать LeaveGroup)
	if userId == requesterId {
		return errors.New("use LeaveGroup to exit the group")
	}

	// Проверяем, не пытается ли админ удалить другого админа
	userRole, err := gr.groupRepo.GetUserRoleInGroup(ctx, groupId, userId)
	if err != nil {
		return err
	}
	if userRole == string(groupMembers.Admin) {
		return errors.New("cannot remove another admin. Change their role to user first")
	}

	// Удаляем пользователя
	err = gr.groupRepo.RemoveUserFromGroup(ctx, groupId, userId)
	if err != nil {
		log.Printf("RemoveUserFromGroup | Failed to remove user from group: %v", err)
		return err
	}

	return nil
}

func (gr *group) CheckUserInGroup(ctx context.Context, groupId, userId int) (bool, error) {
	if groupId == 0 {
		return false, errors.New("group id cannot be 0")
	}
	if userId == 0 {
		return false, errors.New("user id cannot be 0")
	}

	exists, err := gr.groupRepo.CheckUserInGroup(ctx, groupId, userId)
	if err != nil {
		log.Printf("CheckUserInGroup | Failed to check user in group: %v", err)
		return false, err
	}

	return exists, nil
}

func (g *group) AddUserToGroup(ctx context.Context, member groupMembers.GroupMembers, userInfo internal.UserInfo) error {
	if member.GroupId == 0 || member.UserId == 00 || member.Del == 1 {
		return errors.New("data is empty")
	}

	if member.Status != "Admin" {
		member.Status = string(groupMembers.User)
	}

	err := g.groupRepo.AddUserToGroup(member)
	if err != nil {
		return err
	}

	return nil
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

package repository

import (
	"bill-split/internal/domain/entity/group"
	"bill-split/internal/domain/entity/groupMembers"
	"bill-split/internal/domain/entity/user"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepository interface {
	// Транзакция
	TransactionBegin(ctx context.Context) error
	Commit(ctx context.Context) error
	TransactionRollback(ctx context.Context) error

	// Изменить данные (добавить, удалить, обновить)
	CreateGroup(groupData group.Group) (int, error)
	AddUserToGroup(groupData groupMembers.GroupMembers) error
	RemoveUserFromGroup(ctx context.Context, groupId, userId int) error

	// Получить данные
	GetGroupIdByLink(link string) (int, error)
	GetGroup(groupId int) (group.Group, error)
	GetGroupsByUserId(ctx context.Context, userId int) ([]group.Group, error)
	GetMembersByGroupId(ctx context.Context, groupId int) ([]groupMembers.GroupMembers, error)
	CheckUserInGroup(ctx context.Context, groupId, userId int) (bool, error)
	GetUserRoleInGroup(ctx context.Context, groupId, userId int) (string, error)
}

type groupRepository struct {
	db   *pgxpool.Pool
	dbTx pgx.Tx
}

func NewGroupRepository(db *pgxpool.Pool) GroupRepository {
	return &groupRepository{
		db: db,
	}
}

func (gr *groupRepository) GetGroupIdByLink(link string) (int, error) {
	var id int
	var err error
	query := `SELECT id FROM "group" WHERE link_invite = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if gr.dbTx != nil {
		err = gr.dbTx.QueryRow(
			ctx,
			query,
			link,
		).Scan(
			&id,
		)
	} else {
		err = gr.db.QueryRow(
			ctx,
			query,
			link,
		).Scan(
			&id,
		)
	}

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (gr *groupRepository) GetUserRoleInGroup(ctx context.Context, groupId, userId int) (string, error) {
	var status string
	query := `
		SELECT status 
		FROM group_members 
		WHERE group_id = $1 
		  AND user_id = $2 
		  AND del = 0
	`

	if gr.dbTx != nil {
		err := gr.dbTx.QueryRow(ctx, query, groupId, userId).Scan(&status)
		if err != nil {
			return "", err
		}
	} else {
		err := gr.db.QueryRow(ctx, query, groupId, userId).Scan(&status)
		if err != nil {
			return "", err
		}
	}

	return status, nil
}

func (gr *groupRepository) RemoveUserFromGroup(ctx context.Context, groupId, userId int) error {
	var err error
	query := `
		UPDATE group_members 
		SET del = 1 
		WHERE group_id = $1 
		  AND user_id = $2
	`

	if gr.dbTx != nil {
		_, err = gr.dbTx.Exec(ctx, query, groupId, userId)
	} else {
		_, err = gr.db.Exec(ctx, query, groupId, userId)
	}

	if err != nil {
		return err
	}

	return nil
}

func (gr *groupRepository) CheckUserInGroup(ctx context.Context, groupId, userId int) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM group_members 
			WHERE group_id = $1 
			  AND user_id = $2 
			  AND del = 0
		)
	`

	if gr.dbTx != nil {
		err := gr.dbTx.QueryRow(ctx, query, groupId, userId).Scan(&exists)
		if err != nil {
			return false, err
		}
	} else {
		err := gr.db.QueryRow(ctx, query, groupId, userId).Scan(&exists)
		if err != nil {
			return false, err
		}
	}

	return exists, nil
}

// Методы для транзакции
func (gr *groupRepository) TransactionBegin(ctx context.Context) error {
	var err error
	gr.dbTx, err = gr.db.Begin(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (gr *groupRepository) Commit(ctx context.Context) error {
	if gr.dbTx == nil {
		return errors.New("db transaction is nil")
	}

	err := gr.dbTx.Commit(ctx)
	if err != nil {
		return err
	}

	gr.dbTx = nil
	return nil
}

func (gr *groupRepository) TransactionRollback(ctx context.Context) error {
	if gr.dbTx == nil {
		return errors.New("db transaction is nil")
	}

	err := gr.dbTx.Rollback(ctx)
	if err != nil {
		return err
	}

	gr.dbTx = nil
	return nil
}

// Методы запросов
func (gr *groupRepository) GetGroup(groupId int) (group.Group, error) {
	var groupData group.Group
	var err error
	query := `SELECT id, name, create_at, date_start, date_end FROM "group" WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if gr.dbTx != nil {
		err = gr.dbTx.QueryRow(
			ctx,
			query,
			&groupId,
		).Scan(
			&groupData.Id,
			&groupData.Name,
			&groupData.CreateAt,
			&groupData.DateStart,
			&groupData.DateEnd,
		)
	} else {
		err = gr.db.QueryRow(
			ctx,
			query,
			&groupId,
		).Scan(
			&groupData.Id,
			&groupData.Name,
			&groupData.CreateAt,
			&groupData.DateStart,
			&groupData.DateEnd,
		)
	}

	if err != nil {
		return group.Group{}, err
	}

	return groupData, nil
}

func (gr *groupRepository) CreateGroup(groupData group.Group) (int, error) {
	var err error
	query := `INSERT INTO "group" (name, create_at, date_start, date_end) VALUES ($1, $2, $3, $4) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if gr.dbTx != nil {
		err = gr.dbTx.QueryRow(
			ctx,
			query,
			&groupData.Name,
			&groupData.CreateAt,
			&groupData.DateStart,
			&groupData.DateEnd,
		).Scan(&groupData.Id)
	} else {
		err = gr.db.QueryRow(
			ctx,
			query,
			&groupData.Name,
			&groupData.CreateAt,
			&groupData.DateStart,
			&groupData.DateEnd,
		).Scan(&groupData.Id)
	}

	if err != nil {
		return 0, err
	}

	return groupData.Id, nil
}

func (gr *groupRepository) AddUserToGroup(groupData groupMembers.GroupMembers) error {
	var err error
	query := `INSERT INTO "group_members" (user_id, group_id, money_spent, status) VALUES ($1, $2, $3, $4)`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if gr.dbTx != nil {
		_, err = gr.dbTx.Exec(
			ctx,
			query,
			&groupData.UserId,
			&groupData.GroupId,
			&groupData.Amount,
			&groupData.Status,
		)
	} else {
		_, err = gr.db.Exec(
			ctx,
			query,
			&groupData.UserId,
			&groupData.GroupId,
			&groupData.Amount,
			&groupData.Status,
		)
	}
	if err != nil {
		return err
	}
	return nil
}

func (gr *groupRepository) DelUserFromGroup(userId, groupId int) error {
	var err error
	query := `
		UPDATE "group_members" 
		SET del = 1 
		WHERE user_id = $1 
		  AND group_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if gr.dbTx != nil {
		_, err = gr.dbTx.Exec(
			ctx,
			query,
			userId,
			groupId,
		)
	} else {
		_, err = gr.db.Exec(
			ctx,
			query,
			userId,
			groupId,
		)
	}
	if err != nil {
		return err
	}
	return nil
}

func (gr *groupRepository) GetUsers(groupId int) ([]user.User, error) {
	var err error
	var member user.User
	var members []user.User
	var rows pgx.Rows

	query := `
		SELECT gu.user_id 
		FROM "group_members" gu 
		LEFT JOIN "group" g ON g.id = gu.group_id
		WHERE group_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if gr.dbTx != nil {
		rows, err = gr.dbTx.Query(ctx, query, groupId)
	} else {
		rows, err = gr.db.Query(ctx, query, groupId)
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&member.Id)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

// Получить группы в котором состоит пользователь
func (gr *groupRepository) GetGroupsByUserId(ctx context.Context, userId int) ([]group.Group, error) {
	var err error
	var groupData group.Group
	var groupsData []group.Group
	var rows pgx.Rows

	query := `
		SELECT gu.group_id, g.name, g.create_at, g.date_start, g.date_end
		FROM group_members gu 
		LEFT JOIN group g ON g.id = gu.group_id
		WHERE gu.user_id = $1`

	if gr.dbTx != nil {
		rows, err = gr.dbTx.Query(ctx, query, userId)
	} else {
		rows, err = gr.db.Query(ctx, query, userId)
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(
			&groupData.Id,
			&groupData.Name,
			&groupData.CreateAt,
			&groupData.DateStart,
			&groupData.DateEnd,
		)
		if err != nil {
			return nil, err
		}

		groupsData = append(groupsData, groupData)
	}

	return groupsData, nil
}

// Получить пользователей в группе
func (gr *groupRepository) GetMembersByGroupId(ctx context.Context, groupId int) ([]groupMembers.GroupMembers, error) {
	var err error
	var groupMemberData groupMembers.GroupMembers
	var groupMembersData []groupMembers.GroupMembers
	var rows pgx.Rows

	query := `
		SELECT gu.user_id, gu.group_id, gu.status, gu.money_spent
		FROM group_members gu 
		WHERE gu.group_id = $1 AND gu.del = 0
		`

	if gr.dbTx != nil {
		rows, err = gr.dbTx.Query(ctx, query, groupId)
	} else {
		rows, err = gr.db.Query(ctx, query, groupId)
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(
			&groupMemberData.UserId,
			&groupMemberData.GroupId,
			&groupMemberData.Status,
			&groupMemberData.Amount,
		)
		if err != nil {
			return nil, err
		}

		groupMembersData = append(groupMembersData, groupMemberData)
	}

	return groupMembersData, nil
}

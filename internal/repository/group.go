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
	TransactionBegin(ctx context.Context) error
	Commit(ctx context.Context) error
	TransactionRollback(ctx context.Context) error
	CreateGroup(groupData group.Group) (int, error)
	AddUserToGroup(groupData groupMembers.GroupMembers) error
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
func (gr *groupRepository) CreateGroup(groupData group.Group) (int, error) {
	var err error
	query := `INSERT INTO groups (name, create_at, date_start, date_end) VALUES ($1, $2, $3, $4) RETURNING id`

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
	query := `INSERT INTO group_users (user_id, group_id, money_spent, status) VALUES ($1, $2, $3, $4)`

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
		UPDATE group_users 
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

	defer rows.Close()

	query := `
		SELECT gu user_id 
		FROM group_users gu 
		LEFT JOIN groups g ON g.id = gu.group_id
		WHERE group_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if gr.dbTx != nil {
		rows, err = gr.dbTx.Query(ctx, query, groupId)
	} else {
		rows, err = gr.db.Query(ctx, query, groupId)
	}
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

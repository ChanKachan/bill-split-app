package repository

import (
	"bill-split/internal/domain/entity/cost"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CostRepository interface {
	CreateCost(ctx context.Context, costData cost.Cost) (int, error)
	GetCostByID(ctx context.Context, id int) (*cost.Cost, error)
	GetCostsByGroup(ctx context.Context, groupID int) ([]cost.Cost, error)
	GetCostsByUser(ctx context.Context, userID int) ([]cost.Cost, error)
	UpdateCost(ctx context.Context, costData cost.Cost) error
	DeleteCost(ctx context.Context, id int) error
}

type costRepository struct {
	db *pgxpool.Pool
}

func NewCostRepository(db *pgxpool.Pool) CostRepository {
	return &costRepository{
		db: db,
	}
}

func (r *costRepository) CreateCost(ctx context.Context, costData cost.Cost) (int, error) {
	query := `
		INSERT INTO cost (user_id, group_id, description, sum)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query,
		costData.UserId, costData.GroupId, costData.Description, costData.Sum,
	).Scan(&costData.Id)

	if err != nil {
		return 0, err
	}

	return costData.Id, nil
}

func (r *costRepository) GetCostByID(ctx context.Context, id int) (*cost.Cost, error) {
	query := `
		SELECT id, user_id, group_id, description, sum
		FROM cost
		WHERE id = $1
	`

	var cost cost.Cost
	err := r.db.QueryRow(ctx, query, id).Scan(
		&cost.Id, &cost.UserId, &cost.GroupId, &cost.Description, &cost.Sum,
	)

	if err != nil {
		return nil, err
	}

	return &cost, nil
}

func (r *costRepository) GetCostsByGroup(ctx context.Context, groupID int) ([]cost.Cost, error) {
	query := `
		SELECT id, user_id, group_id, description, sum
		FROM cost
		WHERE group_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var costs []cost.Cost
	for rows.Next() {
		var c cost.Cost
		err := rows.Scan(&c.Id, &c.UserId, &c.GroupId, &c.Description, &c.Sum)
		if err != nil {
			return nil, err
		}
		costs = append(costs, c)
	}

	return costs, nil
}

func (r *costRepository) GetCostsByUser(ctx context.Context, userID int) ([]cost.Cost, error) {
	query := `
		SELECT id, user_id, group_id, description, sum
		FROM cost
		WHERE user_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var costs []cost.Cost
	for rows.Next() {
		var c cost.Cost
		err := rows.Scan(&c.Id, &c.UserId, &c.GroupId, &c.Description, &c.Sum)
		if err != nil {
			return nil, err
		}
		costs = append(costs, c)
	}

	return costs, nil
}

func (r *costRepository) UpdateCost(ctx context.Context, costData cost.Cost) error {
	if costData.Id == 0 {
		return errors.New("cannot update cost: haven't costId")
	}

	query := `
		UPDATE cost
		SET description = $1, sum = $2
		WHERE id = $3 AND user_id = $4
	`

	result, err := r.db.Exec(ctx, query,
		costData.Description, costData.Sum, costData.Id, costData.UserId,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("cost not found or you don't have permission to update it")
	}

	return nil
}

func (r *costRepository) DeleteCost(ctx context.Context, id int) error {
	query := `DELETE FROM cost WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("cost not found")
	}

	return nil
}

package repository

import (
	"bill-split/internal/domain/entity/user"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUserById(id int) (*user.User, error)
	CreateUser(userData user.User) (int64, error)
	UpdateUser(userData user.User) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) CreateUser(userData user.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*15))
	defer cancel()

	err := u.db.QueryRow(
		ctx,
		`INSERT INTO "user" (name, email, phone, login, password) 
						VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		userData.Name, userData.Email, userData.Phone, userData.Login, userData.Password,
	).Scan(&userData.Id)
	if err != nil {
		return 0, err
	}

	return userData.Id, nil
}

func (u *userRepository) UpdateUser(userData user.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*15))
	defer cancel()

	_, err := u.db.Exec(
		ctx,
		`
		UPDATE "user"
		SET name = $1, email = $2, phone = $3, password = $4 
		WHERE id = $5
		`,
		userData.Name, userData.Email, userData.Phone, userData.Password, &userData.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUserById(id int) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*15))
	defer cancel()

	var user user.User

	err := u.db.QueryRow(
		ctx,
		`
		SELECT name, email, phone, login
		FROM "user"
		WHERE id = $1
	`, id,
	).Scan(
		&user.Name, &user.Email, &user.Phone, &user.Login,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

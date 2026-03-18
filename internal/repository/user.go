package repository

import (
	"bill-split/internal/domain/entity/user"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
	"strings"
	"time"
)

type UserRepository interface {
	GetUserById(id int) (*user.User, error)
	CreateUser(userData user.User) (int, error)
	GetUserIdByLogin(login string) (int, error)
	UpdateUser(ctx context.Context, userData user.User) error
	GetUserByLogin(login string) (*user.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) GetUserByLogin(login string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var user user.User

	err := u.db.QueryRow(ctx, `
		SELECT id, name, email, phone, login, password
		FROM "user"
		WHERE login = $1
	`, login,
	).Scan(
		&user.Id, &user.Name, &user.Email, &user.Phone, &user.Login, &user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) CreateUser(userData user.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := u.db.QueryRow(ctx,
		`INSERT INTO "user" (name, email, phone, login, password) 
						VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		userData.Name, userData.Email, userData.Phone, userData.Login, userData.Password,
	).Scan(&userData.Id)
	if err != nil {
		return 0, err
	}

	return userData.Id, nil
}

// Метод обновляет данные пользователя
func (u *userRepository) UpdateUser(ctx context.Context, userData user.User) error {
	var query strings.Builder
	count := 1
	var args []any
	if userData.Id != 0 {
		return errors.New("cannot update user: haven't userId")
	}

	query.WriteString(`
			UPDATE "user"
		`)
	query.WriteString(`
			SET
		`)

	valueUserElements := reflect.ValueOf(userData)
	typeOfUserElements := reflect.TypeOf(userData)

	// Цикл пробегает по полям в структуре
	for i := 0; i < valueUserElements.NumField(); i++ {
		if typeOfUserElements.Field(i).Type.Kind() == reflect.String {
			if valueUserElements.Field(i).String() != "" {
				structTag := typeOfUserElements.Field(i).Tag.Get("db")
				if structTag != "" && structTag != "-" {
					if count != 1 {
						query.WriteString(",")
					}
					query.WriteString(fmt.Sprintf("%s = $%d", structTag, count))
					args = append(args, valueUserElements.Field(i).String())
					count++
				}
			}
		}
	}

	if len(args) == 0 {
		return nil
	}

	query.WriteString(fmt.Sprintf("WHERE id = $%d", count))
	args = append(args, userData.Id)
	count++

	_, err := u.db.Exec(
		ctx,
		query.String(),
		args...,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUserById(id int) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var user user.User

	err := u.db.QueryRow(ctx, `
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

func (u *userRepository) GetUserIdByLogin(login string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var id int

	err := u.db.QueryRow(ctx, `
		SELECT id
		FROM "user"
		WHERE login = $1
	`, login,
	).Scan(
		&id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
}

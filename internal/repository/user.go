package repository

import (
	"bill-split/internal/domain/entity/user"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetUserById(id int) (*user.User, error)
	CreateUser(userData user.User) (int, error)
	GetUserIdByLogin(login string) (int, error)
	UpdateUser(userData user.User) error
	GetUserByLogin(login string) (*user.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) GetUserByLogin(login string) (*user.User, error) {
	var user user.User

	err := u.db.QueryRow(`
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
	err := u.db.QueryRow(
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
	_, err := u.db.Exec(`
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
	var user user.User

	err := u.db.QueryRow(`
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
	var id int

	err := u.db.QueryRow(`
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

package service

import (
	"bill-split/internal/domain/entity/user"
	"bill-split/internal/repository"
	"context"
)

type UserHttpService interface {
	UpdateUser(ctx context.Context, userData user.User) error
	GetUser(ctx context.Context, userID int) (*user.User, error)
}
type userHttpService struct {
	userRepo repository.UserRepository
}

func NewUserHttpService(userRepo repository.UserRepository) UserHttpService {
	return &userHttpService{
		userRepo: userRepo,
	}
}

func (u *userHttpService) UpdateUser(ctx context.Context, userData user.User) error {
	err := u.userRepo.UpdateUser(ctx, userData)
	if err != nil {
		return err
	}

	return nil
}

func (u *userHttpService) GetUser(ctx context.Context, userID int) (*user.User, error) {
	userData, err := u.userRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

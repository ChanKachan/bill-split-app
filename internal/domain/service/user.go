package service

import (
	"bill-split/internal/domain/entity/user"
	"bill-split/internal/repository"
	"bill-split/internal/utils"
	"context"
)

type UserHttpService interface {
	UpdateUser(ctx context.Context, userData user.User) error
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
	passwordHash, err := utils.HashPassword(userData.Password)
	if err != nil {
		return err
	}

	userData.Password = passwordHash

	err = u.userRepo.UpdateUser(ctx, userData)
	if err != nil {
		return err
	}

	return nil
}

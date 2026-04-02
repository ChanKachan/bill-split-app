package service

import (
	"bill-split/internal/domain/entity/cost"
	"bill-split/internal/repository"
	"context"
	"errors"
)

type CostService interface {
	CreateCost(ctx context.Context, costData cost.Cost) (int, error)
	GetCost(ctx context.Context, id int) (*cost.Cost, error)
	GetGroupCosts(ctx context.Context, groupID int) ([]cost.Cost, error)
	GetUserCosts(ctx context.Context, userID int) ([]cost.Cost, error)
	UpdateCost(ctx context.Context, costData cost.Cost) error
	DeleteCost(ctx context.Context, id int) error
}

type costService struct {
	costRepo repository.CostRepository
	groupSvc GroupService
}

func NewCostService(costRepo repository.CostRepository, groupSvc GroupService) CostService {
	return &costService{
		costRepo: costRepo,
		groupSvc: groupSvc,
	}
}

func (s *costService) CreateCost(ctx context.Context, costData cost.Cost) (int, error) {
	// Проверяем, что пользователь состоит в группе
	isMember, err := s.groupSvc.CheckUserInGroup(ctx, costData.GroupId, costData.UserId)
	if err != nil {
		return 0, err
	}
	if !isMember {
		return 0, errors.New("user is not a member of this group")
	}

	return s.costRepo.CreateCost(ctx, costData)
}

func (s *costService) GetCost(ctx context.Context, id int) (*cost.Cost, error) {
	return s.costRepo.GetCostByID(ctx, id)
}

func (s *costService) GetGroupCosts(ctx context.Context, groupID int) ([]cost.Cost, error) {
	return s.costRepo.GetCostsByGroup(ctx, groupID)
}

func (s *costService) GetUserCosts(ctx context.Context, userID int) ([]cost.Cost, error) {
	return s.costRepo.GetCostsByUser(ctx, userID)
}

func (s *costService) UpdateCost(ctx context.Context, costData cost.Cost) error {
	// Проверяем, что пользователь является создателем расхода
	existingCost, err := s.costRepo.GetCostByID(ctx, costData.Id)
	if err != nil {
		return err
	}
	if existingCost.UserId != costData.UserId {
		return errors.New("you can only update your own costs")
	}

	return s.costRepo.UpdateCost(ctx, costData)
}

func (s *costService) DeleteCost(ctx context.Context, id int) error {
	return s.costRepo.DeleteCost(ctx, id)
}

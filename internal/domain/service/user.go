package service

import (
	"bill-split/internal/domain/entity/user"
	"bill-split/internal/repository"
	"bill-split/internal/utils"
	proto "bill-split/proto/this"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService struct {
	userRepo repository.UserRepository
	proto.UnimplementedUserServiceServer
}

func NewUserService(userRepo repository.UserRepository) proto.UserServiceServer {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	if req.Data.Email == "" || req.Data.Phone == "" || req.Data.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "userService | empty user data")
	}

	passwordHash, err := utils.HashPassword(req.Data.Password)
	if err != nil {
		return nil, err
	}

	var user = user.User{
		Name:     req.Data.Name,
		Email:    req.Data.Email,
		Phone:    req.Data.Phone,
		Login:    req.Data.Login,
		Password: passwordHash,
	}

	userId, err := u.userRepo.CreateUser(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.CreateUserResponse{
		Id:   userId,
		Code: 200,
	}, nil
}

func (u *userService) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	if req.Data.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "userService | empty user data")
	}

	passwordHash, err := utils.HashPassword(req.Data.Password)
	if err != nil {
		return nil, err
	}

	var user = user.User{
		Id:       req.Data.Id,
		Name:     req.Data.Name,
		Email:    req.Data.Email,
		Phone:    req.Data.Phone,
		Password: passwordHash,
	}

	err = u.userRepo.UpdateUser(user)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "userService | GetUserById error: %v", err)
	}

	return &proto.UpdateUserResponse{Code: 200}, nil
}

func (u *userService) GetUserById(ctx context.Context, request *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	if request.Data.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "userService | empty user id")
	}

	user, err := u.userRepo.GetUserById(int(request.Data.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "userService | GetUserById error: %v", err)
	}

	return &proto.GetUserResponse{
		Name:  user.Name,
		Phone: user.Phone,
		Email: user.Email,
		Login: user.Login,
		Code:  200,
	}, nil
}

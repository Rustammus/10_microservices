package grps

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/task_5/server/internal/models"
	"microservices/task_5/server/internal/services"
	"microservices/task_5/server/proto/user"
)

type userService struct {
	user.UnimplementedUserServer
	s *services.Services
}

func (s *userService) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	claims, err := getClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if req.Id != claims.UserID {
		return nil, errPermissionDenied
	}

	userFound, err := s.s.User.FindById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &user.GetUserResponse{
		Id:    userFound.ID,
		Name:  userFound.Name,
		Email: userFound.Email,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	claims, err := getClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if req.Id != claims.UserID {
		return nil, errPermissionDenied
	}

	err = s.s.User.Update(ctx, models.User{
		ID:       req.Id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &user.UpdateUserResponse{Success: true}, nil
}

func (s *userService) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	claims, err := getClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if req.Id != claims.UserID {
		return nil, errPermissionDenied
	}

	err = s.s.User.DeleteById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &user.DeleteUserResponse{Success: true}, nil
}

func getClaims(ctx context.Context) (models.Claims, error) {
	claimsRaw := ctx.Value(ClaimsKey)
	if claimsRaw == nil {
		return models.Claims{}, errors.New("no claims found in ctx")
	}
	claims, ok := claimsRaw.(models.Claims)
	if !ok {
		return models.Claims{}, errors.New("invalid type found in ctx by ClaimsKey")
	}
	return claims, nil
}

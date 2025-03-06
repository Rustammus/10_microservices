package grps

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/internal/services"
	"awesomeProject1/internal/services/user"
	"awesomeProject1/proto/auth"
	"context"
	"database/sql"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authService struct {
	auth.UnimplementedAuthServer
	s *services.Services
}

func (s *authService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	id, err := s.s.User.Register(ctx, models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, user.ErrEmailUsed) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.RegisterResponse{Id: id}, nil
}

func (s *authService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	userFound, err := s.s.User.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.Unauthenticated, "Incorrect password or email")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if userFound.Password != req.Password {
		return nil, status.Error(codes.Unauthenticated, "Incorrect password or email")
	}

	token, err := s.s.Auth.GetTokens(ctx, userFound)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.LoginResponse{Token: token.AccessToken}, nil
}

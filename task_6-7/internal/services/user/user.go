package user

import (
	"context"
	"database/sql"
	"errors"
	"microservices/task_6/internal/models"
	"microservices/task_6/internal/repository"
)

type Service struct {
	repo repository.UserRepo
}

func NewService(repo repository.UserRepo) *Service {
	return &Service{repo: repo}
}

var ErrEmailUsed = errors.New("email is used")

func (s *Service) Register(ctx context.Context, user models.User) (int64, error) {
	_, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.repo.Create(ctx, user)
		}
		return 0, err
	}

	return 0, ErrEmailUsed
}

func (s *Service) FindById(ctx context.Context, id int64) (models.User, error) {
	return s.repo.Read(ctx, id)
}

func (s *Service) FindByEmail(ctx context.Context, email string) (models.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *Service) Update(ctx context.Context, user models.User) error {
	return s.repo.Update(ctx, user)
}

func (s *Service) DeleteById(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

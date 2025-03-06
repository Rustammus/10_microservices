package user

import (
	"context"
	"database/sql"
	"errors"
	"microservices/task_5/server/internal/models"
	"microservices/task_5/server/internal/repository"
	"microservices/task_5/server/internal/services/notify"
)

type Service struct {
	repo   repository.UserRepo
	notify *notify.Service
}

func NewService(repo repository.UserRepo, notify *notify.Service) *Service {
	return &Service{repo: repo, notify: notify}
}

var ErrEmailUsed = errors.New("email is used")

func (s *Service) Register(ctx context.Context, user models.User) (int64, error) {
	_, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			id, err := s.repo.Create(ctx, user)
			user.ID = id
			s.notify.UserCreated(user)
			return id, err
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
	userOld, userNew, err := s.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	s.notify.UserUpdated(userOld, userNew)
	return nil
}

func (s *Service) DeleteById(ctx context.Context, id int64) error {
	delUser, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	s.notify.UserDeleted(delUser)
	return nil
}

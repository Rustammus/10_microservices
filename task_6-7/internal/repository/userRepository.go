package repository

import (
	"context"
	"microservices/task_6/internal/models"
)

type UserRepo interface {
	Create(ctx context.Context, user models.User) (int64, error)
	Read(ctx context.Context, id int64) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id int64) error
}

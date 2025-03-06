package repository

import (
	"context"
	"microservices/task_5/server/internal/models"
)

type UserRepo interface {
	Create(ctx context.Context, user models.User) (int64, error)
	Read(ctx context.Context, id int64) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	Update(ctx context.Context, user models.User) (userOld, userNew models.User, err error)
	Delete(ctx context.Context, id int64) (deleted models.User, err error)
}

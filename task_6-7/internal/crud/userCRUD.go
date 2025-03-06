package crud

import (
	"context"
	"database/sql"
	"log/slog"
	"math/rand"
	"microservices/task_6/internal/models"
)

type UserCRUD struct {
	m map[int64]models.User
	l *slog.Logger
}

func NewUserCRUD(l *slog.Logger) *UserCRUD {
	l = l.With(slog.String("component", "internal/service/userCRUD"))

	return &UserCRUD{
		m: make(map[int64]models.User),
		l: l,
	}
}

func (c *UserCRUD) Create(ctx context.Context, user models.User) (int64, error) {
	id := rand.Int63()
	user.ID = id
	c.m[id] = user

	c.l.DebugContext(ctx, "created user", slog.Any("user", user))

	return id, nil
}

func (c *UserCRUD) Read(ctx context.Context, id int64) (models.User, error) {
	u, ok := c.m[id]
	if !ok {
		return models.User{}, sql.ErrNoRows
	}

	c.l.DebugContext(ctx, "read user", slog.Any("user", u))

	return u, nil
}

func (c *UserCRUD) FindByEmail(ctx context.Context, email string) (models.User, error) {
	for _, u := range c.m {
		if u.Email == email {
			c.l.DebugContext(ctx, "found user by email", slog.Any("user", u))
			return u, nil
		}
	}
	return models.User{}, sql.ErrNoRows
}

func (c *UserCRUD) Update(ctx context.Context, user models.User) error {
	userOld, ok := c.m[user.ID]
	if !ok {
		return sql.ErrNoRows
	}

	userNew := models.User{
		ID:       userOld.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	c.m[user.ID] = userNew

	c.l.DebugContext(ctx, "updated user", slog.Any("user_old", userOld), slog.Any("user_new", userNew))
	return nil
}

func (c *UserCRUD) Delete(ctx context.Context, id int64) error {
	user, ok := c.m[id]
	if !ok {
		return sql.ErrNoRows
	}
	delete(c.m, id)

	c.l.DebugContext(ctx, "deleted user", slog.Any("user", user))

	return nil
}

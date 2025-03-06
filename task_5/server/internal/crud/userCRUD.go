package crud

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"microservices/task_5/server/internal/models"
)

type UserCRUD struct {
	m map[int64]models.User
}

func NewUserCRUD() *UserCRUD {
	return &UserCRUD{
		m: make(map[int64]models.User),
	}
}

func (c *UserCRUD) Create(_ context.Context, user models.User) (int64, error) {
	id := rand.Int63()
	user.ID = id
	c.m[id] = user

	log.Print("create user", user)
	return id, nil
}

func (c *UserCRUD) Read(_ context.Context, id int64) (models.User, error) {
	u, ok := c.m[id]
	if !ok {
		return models.User{}, sql.ErrNoRows
	}

	log.Print("read user", u)
	return u, nil
}

func (c *UserCRUD) FindByEmail(_ context.Context, email string) (models.User, error) {
	for _, u := range c.m {
		if u.Email == email {
			log.Print("found user by email ", u)
			return u, nil
		}
	}
	return models.User{}, sql.ErrNoRows
}

func (c *UserCRUD) Update(_ context.Context, user models.User) (userOld, userNew models.User, err error) {
	userOld, ok := c.m[user.ID]
	if !ok {
		return models.User{}, models.User{}, sql.ErrNoRows
	}

	userNew = models.User{
		ID:       userOld.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	c.m[user.ID] = userNew
	log.Print("update user ", userNew)
	return userOld, userNew, nil
}

func (c *UserCRUD) Delete(_ context.Context, id int64) (deleted models.User, err error) {
	deleted, ok := c.m[id]
	if !ok {
		return models.User{}, sql.ErrNoRows
	}
	delete(c.m, id)
	log.Print("delete user ", id)

	return deleted, nil
}

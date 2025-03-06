package services

import (
	"microservices/task_6/internal/services/auth"
	"microservices/task_6/internal/services/user"
)

type Services struct {
	User *user.Service
	Auth *auth.Service
}

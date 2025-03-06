package services

import (
	"microservices/task_5/server/internal/services/auth"
	"microservices/task_5/server/internal/services/user"
)

type Services struct {
	User *user.Service
	Auth *auth.Service
}

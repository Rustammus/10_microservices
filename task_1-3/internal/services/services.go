package services

import (
	"awesomeProject1/internal/services/auth"
	"awesomeProject1/internal/services/user"
)

type Services struct {
	User *user.Service
	Auth *auth.Service
}

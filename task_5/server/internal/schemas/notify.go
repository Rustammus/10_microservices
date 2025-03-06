package schemas

import "microservices/task_5/server/internal/models"

type NotifyConfirm struct {
	Id int64 `json:"id"`
}

type NotifySend struct {
	Id          int64              `json:"id"`
	UserCreated *NotifyUserCreated `json:"user_created"`
	UserUpdated *NotifyUserUpdated `json:"user_updated"`
	UserDeleted *NotifyUserDeleted `json:"user_deleted"`
}

type NotifyUserCreated struct {
	User models.User `json:"user"`
}

type NotifyUserUpdated struct {
	UserOld models.User `json:"user_old"`
	UserNew models.User `json:"user_new"`
}

type NotifyUserDeleted struct {
	User models.User `json:"user"`
}

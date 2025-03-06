package models

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

//func (u *User) LogValue() slog.Value {
//	//TODO implement me
//	panic("implement me")
//}

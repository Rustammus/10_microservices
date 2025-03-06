package schemas

type UserRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	ID int64 `json:"id"`
}

type UserFindByIdResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

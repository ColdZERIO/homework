package handler

type PersistUserRequest struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type UserResponse struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt int64 `json:"created_at"`
}

type UserListRequest struct {
	Users  []UserResponse `json:"users"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

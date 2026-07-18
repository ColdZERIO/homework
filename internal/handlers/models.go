package handler

import "time"

type PersistUserRequest struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type PersistUserResponse struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
}

type UserListRequest struct {
	Users  []UserResponse `json:"users"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

type FindUserResponse struct {
	ID        string    `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type FindUserRequest struct {
	ID        string    `json:"id"`
}

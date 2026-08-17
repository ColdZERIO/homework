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
	ID    string `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type FindUserRequest struct {
	ID string `json:"id"`
}

type FindUserResponse struct {
	ID        string    `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type FindUserListRequest struct {
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Where   string `json:"where"`
	OrderBy string `json:"orderby"`
}

type FindUserListResponse struct {
	Users []FindUserResponse `json:"users"`
}

type DeleteUserRequest struct {
	ID string `json:"id"`
}

type DeleteUserResponse struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

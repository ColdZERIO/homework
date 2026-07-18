package services

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Login     string
	Name      string
	Email     string
	CreatedAt time.Time
}

type CreateUserInput struct {
	ID       string
	Login    string
	Password string
	Name     string
	Email    string
}

type GetUser struct {
	ID    uuid.UUID
	Login string
	Name  string
	Email string
}

type GetUserList struct {
	Users  []GetUser
	Limit  int
	Offset int
}

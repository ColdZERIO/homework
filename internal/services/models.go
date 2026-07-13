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

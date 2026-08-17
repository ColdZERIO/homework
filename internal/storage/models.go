package storage

import (
	"time"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type UserModel struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Login        string    `gorm:"column:login"`
	PasswordHash string    `gorm:"column:password"`
	Name         string    `gorm:"column:name"`
	Role         string    `gorm:"column:role;default:user"`
	Email        string    `gorm:"column:email"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	IsActive     bool      `gorm:"column:is_active"`
}

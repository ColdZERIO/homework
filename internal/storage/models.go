package storage

import (
	"time"
)

type UserModel struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Login        string    `gorm:"column:login"`
	Role         string    `gorm:"column:role"`
	PasswordHash string    `gorm:"column:password"`
	Name         string    `gorm:"column:name"`
	Email        string    `gorm:"column:email"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	IsActive     bool      `gorm:"column:is_active"`
}

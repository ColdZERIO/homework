package storage

import "github.com/google/uuid"

type UserDB struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Login        string    `gorm:"column:login"`
	PasswordHash string    `gorm:"column:password"`
	Name         string    `gorm:"column:name"`
	Email        string    `gorm:"column:email"`
	CreatedAt    int64     `gorm:"column:created_at"`
	IsActive     bool      `gorm:"column:is_active"`
}

type UserDBResponse struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Login     string    `gorm:"column:login"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	CreatedAt int64     `gorm:"column:created_at"`
}

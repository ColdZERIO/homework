package storage

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"homework/internal/storage"

	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStorage interface {
	Persist(ctx context.Context, userDB storage.UserModel) (UserModel, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Find(ctx context.Context, id uuid.UUID) (storage.UserModel, error)
	GetList(ctx context.Context, limit, offset int) ([]domain.UserOutput, error)
}

type Storage struct {
	db    *gorm.DB
	cache *MemoryCache
}

func NewUserStorage(db *gorm.DB) *Storage {
	return &Storage{
		db:    db,
		cache: UserMemoryCache(),
	}
}

func (UserModel) TableName() string {
	return "users"
}

func (s *Storage) Persist(ctx context.Context, userModel UserModel) (UserModel, error) {
	err := s.db.WithContext(ctx).Save(&userModel).Error
	if err != nil {
		log.Println(err)
		return UserModel{}, err
	}

	s.cache.Clear()

	return userModel, nil
}

func (s *Storage) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.db.WithContext(ctx).Model(&UserDB{}).Where("id = ?", id).Update("is_active", false).Error
	if err != nil {
		log.Println(err)
		return err
	}

	s.cache.Clear()

	return nil
}

func (s *Storage) Find(ctx context.Context, id uuid.UUID) (UserModel, error) {
	var userModel UserModel

	key := fmt.Sprintf("userID: %d", id)

	if value, ok := s.cache.Get(key); ok {
		user := value.(UserModel)
		return user, nil
	}

	err := s.db.WithContext(ctx).First(&userModel, id).Error
	if err != nil {
		log.Println(err)
		return userModel, err
	}

	s.cache.Set(key, userModel)

	return userModel, nil
}

func (s *Storage) GetList(ctx context.Context, limit, offset int) ([]domain.UserOutput, error) {
	key := "users:list"

	if value, ok := s.cache.Get(key); ok {
		users := value.([]domain.UserOutput)
		return users, nil
	}

	var usersDB []UserDB

	err := s.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&usersDB).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	s.cache.Set(key, usersDB)

	return ToUserList(usersDB), nil
}

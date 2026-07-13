package storage

import (
	"context"
	"fmt"

	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

func (UserDB) TableName() string {
	return "users"
}

func (s *Storage) Persist(ctx context.Context, userDB UserDB) (uuid.UUID, error) {
	if userDB.ID == uuid.Nil {
		err := s.db.WithContext(ctx).Create(&userDB).Error
		if err != nil {
			log.Println(err)
			return uuid.Nil, err
		}

		return userDB.ID, nil
	}

	err := s.db.WithContext(ctx).Save(&userDB).Error
	if err != nil {
		log.Println(err)
		return uuid.Nil, err
	}

	s.cache.Clear()

	return userDB.ID, nil
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

func (s *Storage) Find(ctx context.Context, id uuid.UUID) (UserDBResponse, error) {
	key := fmt.Sprintf("userID: %d", id)

	if value, ok := s.cache.Get(key); ok {
		user := value.(UserDBResponse)
		return user, nil
	}

	var userDB UserDB

	err := s.db.WithContext(ctx).First(&userDB, id).Error
	if err != nil {
		log.Println(err)
		return UserDBResponse{}, err
	}

	s.cache.Set(key, ToUser(userDB))

	return ToUser(userDB), nil
}

func (s *Storage) GetList(ctx context.Context, limit, offset int) ([]UserDBResponse, error) {
	key := "users:list"

	if value, ok := s.cache.Get(key); ok {
		users := value.([]UserDBResponse)
		return users, nil
	}

	var usersDB []UserDB

	err := s.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&usersDB).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	users := ToUserList(usersDB)

	s.cache.Set(key, users)

	return users, nil
}

package storage

import (
	"context"
	"fmt"

	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStorage interface {
	Persist(ctx context.Context, userDB UserModel) (UserModel, error)
	Delete(ctx context.Context, id string) (UserModel, error)
	Find(ctx context.Context, id string) (UserModel, error)
	GetList(ctx context.Context, limit, offset int, where, orderby string) ([]UserModel, error)
	AuthUser(ctx context.Context, login string) (UserModel, error)
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
	if userModel.ID == "" {
		userModel.ID = uuid.New().String()
		err := s.db.WithContext(ctx).Create(&userModel).Error
		if err != nil {
			log.Println(err)
			return UserModel{}, err
		}
	}

	err := s.db.WithContext(ctx).Model(&userModel).Where("id = ?", userModel.ID).Updates(userModel).Error
	if err != nil {
		log.Println(err)
		return UserModel{}, err
	}

	s.cache.Set(userIDkey(userModel.ID), userModel)
	s.cache.Set(userLoginKey(userModel.Login), userModel)

	return userModel, nil
}

func (s *Storage) Delete(ctx context.Context, id string) (UserModel, error) {
	var userModel UserModel

	err := s.db.WithContext(ctx).Model(&userModel).Where("id = ?", id).Update("is_active", false).Error
	if err != nil {
		log.Println(err)
		return UserModel{}, err
	}

	s.cache.Delete(userIDkey(userModel.ID))

	return userModel, nil
}

func (s *Storage) Find(ctx context.Context, id string) (UserModel, error) {
	var userModel UserModel

	key := fmt.Sprintf("userID: %s", id)

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

func (s *Storage) AuthUser(ctx context.Context, login string) (UserModel, error) {
	var userModel UserModel

	key := fmt.Sprintf("login: %s", login)

	if value, ok := s.cache.Get(key); ok {
		user := value.(UserModel)
		return user, nil
	}

	err := s.db.WithContext(ctx).First(&userModel, login).Error
	if err != nil {
		log.Println(err)
		return userModel, err
	}

	s.cache.Set(key, userModel)

	return userModel, nil
}

func (s *Storage) GetList(ctx context.Context, limit, offset int, where, orderby string) ([]UserModel, error) {
	key := userListKey(limit, offset, where, orderby)

	if value, ok := s.cache.Get(key); ok {
		users := value.([]UserModel)
		return users, nil
	}

	var userModel []UserModel

	err := s.db.WithContext(ctx).Limit(limit).Offset(offset).Where(where).Order(orderby).Find(&userModel).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	s.cache.Set(key, userModel)

	return userModel, nil
}

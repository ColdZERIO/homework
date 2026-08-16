package storage

import (
	"context"

	"log"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
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
	cache *cache.Cache
}

func NewUserStorage(db *gorm.DB, cache *cache.Cache) *Storage {
	return &Storage{
		db:    db,
		cache: cache,
	}
}

func (UserModel) TableName() string {
	return "users"
}

func (s *Storage) Persist(ctx context.Context, userModel UserModel) (UserModel, error) {
	key := KeyCacheID(userModel.ID)

	if userModel.ID == "" {
		userModel.ID = uuid.New().String()
		err := s.db.WithContext(ctx).Create(&userModel).Error
		if err != nil {
			// log

			return UserModel{}, err
		}

		s.cache.Set(key, userModel, ttl)
	}

	err := s.db.WithContext(ctx).Model(&userModel).Where("id = ?", userModel.ID).Updates(userModel).Error
	if err != nil {
		// log
		return UserModel{}, err
	}

	s.cache.Set(key, userModel, ttl)

	return userModel, nil
}

func (s *Storage) Delete(ctx context.Context, id string) (UserModel, error) {
	var userModel UserModel
	key := KeyCacheID(id)

	err := s.db.WithContext(ctx).Model(&userModel).Where("id = ?", id).Update("is_active", false).Error
	if err != nil {
		// log
		return UserModel{}, err
	}

	s.cache.Delete(key)

	return userModel, nil
}

func (s *Storage) Find(ctx context.Context, id string) (UserModel, error) {
	var userModel UserModel

	key := KeyCacheID(id)
	value, found := s.cache.Get(key)
	if found {
		user, ok := value.(UserModel)
		if ok {
			return user, nil
		}

		s.cache.Delete(key)
	}

	err := s.db.WithContext(ctx).First(&userModel, id).Error
	if err != nil {
		log.Println(err)
		return UserModel{}, err
	}

	s.cache.Set(key, userModel, ttl)

	return userModel, nil
}

func (s *Storage) AuthUser(ctx context.Context, login string) (UserModel, error) { // ???
	var userModel UserModel

	err := s.db.WithContext(ctx).First(&userModel, login).Error
	if err != nil {
		log.Println(err)
		return userModel, err
	}

	return userModel, nil
}

func (s *Storage) GetList(ctx context.Context, limit, offset int, where, orderby string) ([]UserModel, error) {
	var userModel []UserModel
	key := KeyCacheList(QueryParams{
		limit:   limit,
		offset:  offset,
		where:   where,
		orderby: orderby,
	})

	value, found := s.cache.Get(key)
	if found {
		return value.([]UserModel), nil
	}

	err := s.db.WithContext(ctx).Limit(limit).Offset(offset).Where(where).Order(orderby).Find(&userModel).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	s.cache.Set(key, userModel, ttl)

	return userModel, nil
}

package storage

import (
	"context"
	"fmt"
	handler "homework/internal/handlers"
	"homework/internal/model"
	"log"

	"gorm.io/gorm"
)

type UserDB struct {
	ID        int    `gorm:"column:id;primaryKey"`
	Login     string `gorm:"column:login"`
	Password  string `gorm:"column:password"`
	Name      string `gorm:"column:name"`
	role      string `gorm:"column:role"`
	Email     string `gorm:"column:email"`
	CreatedAt int64  `gorm:"column:created_at"`
	IsActive  bool   `gorm:"column:is_active"`
}

type Storage struct {
	db    *gorm.DB
	cache *MemoryCache
}

func UserStorage(db *gorm.DB) *Storage {
	return &Storage{
		db:    db,
		cache: UserMemoryCache(),
	}
}

func (UserDB) TableName() string {
	return "users"
}

func (s *Storage) Persist(ctx context.Context, userDB UserDB) (int, error) {
	err := s.db.WithContext(ctx).Create(&userDB).Error
	if err != nil {
		log.Println(err)
		return 0, err
	}

	s.cache.Clear()

	return userDB.ID, nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	err := s.db.WithContext(ctx).Model(&UserDB{}).Where("id = ?", id).Update("is_active", false).Error
	if err != nil {
		log.Println(err)
		return err
	}

	s.cache.Clear()

	return nil
}

func (s *Storage) Find(ctx context.Context, id int) (model.User, error) {
	key := fmt.Sprintf("userID: %d", id)

	if value, ok := s.cache.Get(key); ok {
		user := value.(model.User)
		return user, nil
	}

	var userDB UserDB
	err := s.db.WithContext(ctx).First(&userDB, id).Error
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	s.cache.Set(key, ToUser(userDB))

	return ToUser(userDB), nil
}

func (s *Storage) Update(ctx context.Context, userReq handler.UserRequest) error {
	userDB := UserDB{
		Name:  userReq.Name,
		Email: userReq.Email,
	}

	err := s.db.WithContext(ctx).Model(&UserDB{}).Where("id = ?", userDB.ID).Updates(UserDB{Name: userDB.Name, Email: userDB.Email}).Error
	if err != nil {
		log.Println(err)
		return err
	}

	s.cache.Clear()

	return nil
}

func (s *Storage) GetList(ctx context.Context, limit, offset int) ([]model.User, error) {
	key := "users:list"

	if value, ok := s.cache.Get(key); ok {
		users := value.([]model.User)
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

func ToUser(userDB UserDB) model.User {
	return model.User{
		ID:        userDB.ID,
		Login:     userDB.Login,
		Name:      userDB.Name,
		Email:     userDB.Email,
		CreatedAt: userDB.CreatedAt,
	}
}

func ToUserList(userListDB []UserDB) []model.User {
	userList := make([]model.User, len(userListDB))

	for _, userDB := range userListDB {
		userList = append(userList, model.User{
			ID:        userDB.ID,
			Login:     userDB.Login,
			Name:      userDB.Name,
			Email:     userDB.Email,
			CreatedAt: userDB.CreatedAt,
		})
	}

	return userList
}

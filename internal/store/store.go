package store

import (
	"context"
	"homework/internal/model"
	"log"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(ctx context.Context, user model.User) (int, error) {
	err := s.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return user.ID, nil
}

func (s *Store) DeleteUser(ctx context.Context, id int) error {
	err := s.db.WithContext(ctx).Delete(&model.User{}, id).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Store) GetUser(ctx context.Context, id int) (model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	return user, nil
}

func (s *Store) UpdateUser(ctx context.Context, user model.User) error {
	err := s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(model.User{Name: user.Name, Email: user.Email}).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Store) GetUsersList(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := s.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

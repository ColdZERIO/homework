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

func (s *Store) CreateUser(ctx context.Context, userDB model.UserDB) (int, error) {
	err := s.db.WithContext(ctx).Create(&userDB).Error
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return userDB.ID, nil
}

func (s *Store) DeleteUser(ctx context.Context, id int) error {
	err := s.db.WithContext(ctx).Delete(&model.UserDB{}, id).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Store) GetUser(ctx context.Context, id int) (model.User, error) {
	var userDB model.UserDB
	err := s.db.WithContext(ctx).First(&userDB, id).Error
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	return ToUser(userDB), nil
}

func (s *Store) UpdateUser(ctx context.Context, userReq model.UserRequest) error {
	userDB := model.UserDB{
		Name:      userReq.Name,
		Email:     userReq.Email,
	}

	err := s.db.WithContext(ctx).Model(&model.UserDB{}).Where("id = ?", userDB.ID).Updates(model.UserDB{Name: userDB.Name, Email: userDB.Email}).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Store) GetUsersList(ctx context.Context) ([]model.User, error) {
	var usersDB []model.UserDB

	err := s.db.WithContext(ctx).Find(&usersDB).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return ToUserList(usersDB), nil
}

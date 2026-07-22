package services

import (
	"context"
	"homework/internal/storage"
	"log"
)

// Переписать структуры на параметры

type UserService interface {
	Persist(ctx context.Context, ID, login, password, name, email string) (storage.UserModel, error)
	Delete(ctx context.Context, userID string) (storage.UserModel, error)
	Find(ctx context.Context, userID string) (storage.UserModel, error)
	GetList(ctx context.Context, limit, offset int, where, orderby string) ([]storage.UserModel, error)
	CheckPassword(ctx context.Context, login, password string) (storage.UserModel, error)
}

type UserServices struct {
	storage storage.UserStorage
}

func NewUserServices(storage storage.UserStorage) *UserServices {
	return &UserServices{storage: storage}
}

func (s *UserServices) Persist(ctx context.Context, ID, login, password, name, email string) (storage.UserModel, error) {
	if ID == "" {
		// Сгенерировать UUID
	}

	userModel := storage.UserModel{
		ID:           ID,
		Login:        login,
		PasswordHash: HashPassword(password),
		Name:         name,
		Email:        email,
		IsActive:     true,
	}

	user, err := s.storage.Persist(ctx, userModel)
	if err != nil {
		log.Println(err)
		return storage.UserModel{}, err
	}

	return user, nil
}

func (s *UserServices) CheckPassword(ctx context.Context, login, password string) (storage.UserModel, error) {
	user, err := s.storage.AuthUser(ctx, login)
	if err != nil {
		log.Println(err)
		return storage.UserModel{}, err
	}

	err = ValidationPassword(password, user.PasswordHash)
	if err != nil {
		log.Println(err)
		return storage.UserModel{}, err
	}

	return user, nil
}

func (s *UserServices) Delete(ctx context.Context, id string) (storage.UserModel, error) {
	userDeleted, err := s.storage.Delete(ctx, id)
	if err != nil {
		return storage.UserModel{}, err
	}

	return userDeleted, nil
}

func (s *UserServices) Find(ctx context.Context, userID string) (storage.UserModel, error) {
	user, err := s.storage.Find(ctx, userID)
	if err != nil {
		return storage.UserModel{}, err
	}

	return user, nil
}

func (s *UserServices) GetList(ctx context.Context, limit, offset int, where, orderby string) ([]storage.UserModel, error) {
	users, err := s.storage.GetList(ctx, limit, offset, where, orderby)
	if err != nil {
		return nil, err
	}

	return users, nil
}

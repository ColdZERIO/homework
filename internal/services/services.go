package services

import (
	"context"
	"homework/internal/storage"
	"log"

	"github.com/google/uuid"
)

// Переписать структуры на параметры

type UserService interface {
	Persist(ctx context.Context, ID, login, password, name, email string) (storage.UserModel, error)
	Delete(ctx context.Context, id int) error
	Find(ctx context.Context, userID string) (storage.UserModel, error)
	Update(ctx context.Context, userReq PersistUserRequest) error
	GetList(ctx context.Context, limit, offset int) ([]UserResponse, error)
}

type UserServices struct {
	storage storage.UserStorage
}

func NewUserServices(store storage.UserStorage) *UserServices {
	return &UserServices{store: store}
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

func (s *UserServices) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.store.Delete(ctx, uid)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServices) Find(ctx context.Context, userID string) (storage.UserModel, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return storage.UserModel{}, err
	}

	user, err := s.store.Find(ctx, id)
	if err != nil {
		return storage.UserModel{}, err
	}

	return user, nil
}

func (s *UserServices) GetList(ctx context.Context, limit, offset int) ([]storage.UserModel, error) {
	users, err := s.store.GetList(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return ToUserListResponse(users), nil
}

package services

import (
	"context"
	handler "homework/internal/handlers"
	"homework/internal/storage"
	"log"
	"time"

	"github.com/google/uuid"
)

type UserStorage interface {
	Persist(ctx context.Context, userDB storage.UserDB) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Find(ctx context.Context, id uuid.UUID) (storage.UserDBResponse, error)
	GetList(ctx context.Context, limit, offset int) ([]storage.UserDBResponse, error)
}

type UserServices struct {
	store UserStorage
}

func NewUserServices(store UserStorage) *UserServices {
	return &UserServices{store: store}
}

func (s *UserServices) Persist(ctx context.Context, userReq handler.PersistUserRequest) (string, error) {
	userDB := storage.UserDB{
		Login:        userReq.Login,
		PasswordHash: HashPassword(userReq.Password),
		Name:         userReq.Name,
		Email:        userReq.Email,
		CreatedAt:    time.Now().Unix(),
		IsActive:     true,
	}

	id, err := uuid.Parse(userReq.ID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	userDB.ID = id

	id, err = s.store.Persist(ctx, userDB)
	if err != nil {
		return "", err
	}

	return id.String(), nil
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

func (s *UserServices) Find(ctx context.Context, userID string) (handler.UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return handler.UserResponse{}, err
	}

	user, err := s.store.Find(ctx, id)
	if err != nil {
		return handler.UserResponse{}, err
	}

	return ToUserResponse(user), nil
}

func (s *UserServices) GetList(ctx context.Context, limit, offset int) ([]handler.UserListRequest, error) {
	users, err := s.store.GetList(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}

package services

import (
	"context"
	handler "homework/internal/handlers"
	"homework/internal/model"
	"homework/internal/storage"
	"strconv"
	"time"
)

type Storage interface {
	Persist(ctx context.Context, userDB storage.UserDB) (int, error)
	Delete(ctx context.Context, id int) error
	FindByID(ctx context.Context, id int) (model.User, error)
	Update(ctx context.Context, userReq handler.UserRequest) error
	GetList(ctx context.Context) ([]model.User, error)
}

type Services struct {
	store Storage
}

func UserServices(store Storage) *Services {
	return &Services{store: store}
}

func (s *Services) Persist(ctx context.Context, userReq handler.UserRequest) (int, error) {
	userDB := storage.UserDB{
		Login:     userReq.Login,
		Password:  HashPassword(userReq.Password),
		Name:      userReq.Name,
		Email:     userReq.Email,
		CreatedAt: time.Now().Unix(),
		IsActive:  true,
	}

	id, err := s.store.Persist(ctx, userDB)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Services) Delete(ctx context.Context, id int) error {
	err := s.store.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Services) FindByID(ctx context.Context, userID string) (model.User, error) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return model.User{}, err
	}

	user, err := s.store.FindByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *Services) Update(ctx context.Context, userReq handler.UserRequest) error {
	err := s.store.Update(ctx, userReq)
	if err != nil {
		return err
	}

	return nil
}

func (s *Services) GetList(ctx context.Context) ([]model.User, error) {
	users, err := s.store.GetList(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

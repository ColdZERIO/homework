package services

import (
	"context"
	"homework/internal/model"
	"strconv"
	"time"
)

type Store interface {
	CreateUser(ctx context.Context, userDB model.UserDB) (int, error)
	DeleteUser(ctx context.Context, id int) error
	GetUser(ctx context.Context, id int) (model.User, error)
	UpdateUser(ctx context.Context, userReq model.UserRequest) error
	GetUsersList(ctx context.Context) ([]model.User, error)
}

type Services struct {
	store Store
}

func NewServices(store Store) *Services {
	return &Services{store: store}
}

func (s *Services) CreateUser(ctx context.Context, userReq model.UserRequest) (int, error) {
	userDB := model.UserDB{
		Login:     userReq.Login,
		Password:  HashPassword(userReq.Password),
		Name:      userReq.Name,
		Email:     userReq.Email,
		CreatedAt: time.Now().Unix(),
	}

	id, err := s.store.CreateUser(ctx, userDB)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Services) DeleteUser(ctx context.Context, id int) error {
	err := s.store.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Services) GetUser(ctx context.Context, userID string) (model.User, error) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return model.User{}, err
	}

	user, err := s.store.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *Services) UpdateUser(ctx context.Context, userReq model.UserRequest) error {
	err := s.store.UpdateUser(ctx, userReq)
	if err != nil {
		return err
	}

	return nil
}

func (s *Services) GetUsersList(ctx context.Context) ([]model.User, error) {
	users, err := s.store.GetUsersList(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

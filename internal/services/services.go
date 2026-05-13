package services

import (
	"context"
	"homework/internal/model"
)

type Store interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	DeleteUser(ctx context.Context, id int) error
	GetUser(ctx context.Context, id int) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	GetUsersList(ctx context.Context) ([]model.User, error)
}

type Services struct {
	store Store
}

func NewServices(store Store) *Services {
	return &Services{store: store}
}

func (s *Services) CreateUser(ctx context.Context, user model.User) (int, error) {
	id, err := s.store.CreateUser(ctx, user)
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

func (s *Services) GetUser(ctx context.Context, id int) (model.User, error) {
	user, err := s.store.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *Services) UpdateUser(ctx context.Context, user model.User) error {
	err := s.store.UpdateUser(ctx, user)
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

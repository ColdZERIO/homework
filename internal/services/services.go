package services

import (
	"context"
	"errors"
	"homework/internal/auth"
	"homework/internal/storage"
	"log/slog"
)

var ErrUserInactive = errors.New("user is inactive")

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
	logger  *slog.Logger
}

func NewUserServices(storage storage.UserStorage, logger *slog.Logger) *UserServices {
	return &UserServices{
		storage: storage,
		logger:  logger.With("layer", "services"),
	}
}

func (s *UserServices) Persist(ctx context.Context, ID, login, password, name, email string) (storage.UserModel, error) {
	role := auth.RoleUser
	if ID != "" {
		currentUser, err := s.storage.Find(ctx, ID)
		if err != nil {
			return storage.UserModel{}, err
		}
		role = currentUser.Role
	}

	hashPassword, err := HashPassword(password)
	if err != nil {
		s.logger.Debug(
			"failed to hash password",
			slog.String("user_id", ID),
			slog.Any("error", err),
		)
		return storage.UserModel{}, err
	}

	userModel := storage.UserModel{
		ID:           ID,
		Login:        login,
		PasswordHash: hashPassword,
		Role:         role,
		Name:         name,
		Email:        email,
		IsActive:     true,
	}

	user, err := s.storage.Persist(ctx, userModel)
	if err != nil {
		s.logger.Debug(
			"failed to add user in database",
			slog.String("user_id", ID),
			slog.Any("error", err),
		)
		return storage.UserModel{}, err
	}

	return user, nil
}

func (s *UserServices) CheckPassword(ctx context.Context, login, password string) (storage.UserModel, error) {
	user, err := s.storage.AuthUser(ctx, login)
	if err != nil {
		s.logger.Debug(
			"failed to get password from database",
			slog.String("login", login),
			slog.Any("error", err),
		)

		return storage.UserModel{}, err
	}

	if !user.IsActive {
		return storage.UserModel{}, ErrUserInactive
	}

	err = ValidationPassword(password, user.PasswordHash)
	if err != nil {
		s.logger.Debug(
			"invalid user password",
			slog.String("login", login),
			slog.Any("error", err),
		)

		return storage.UserModel{}, err
	}

	return user, nil
}

func (s *UserServices) Delete(ctx context.Context, id string) (storage.UserModel, error) {
	userDeleted, err := s.storage.Delete(ctx, id)
	if err != nil {
		s.logger.Debug(
			"failed to delete user",
			slog.String("user_id", id),
			slog.Any("error", err),
		)

		return storage.UserModel{}, err
	}

	return userDeleted, nil
}

func (s *UserServices) Find(ctx context.Context, userID string) (storage.UserModel, error) {
	user, err := s.storage.Find(ctx, userID)
	if err != nil {
		s.logger.Debug(
			"failed ti find user",
			slog.String("user_id", userID),
			slog.Any("error", err),
		)

		return storage.UserModel{}, err
	}

	return user, nil
}

func (s *UserServices) GetList(ctx context.Context, limit, offset int, where, orderby string) ([]storage.UserModel, error) {
	users, err := s.storage.GetList(ctx, limit, offset, where, orderby)
	if err != nil {
		s.logger.Debug(
			"invalid user_list query",
			slog.Any("error", err),
		)

		return nil, err
	}

	return users, nil
}

package services

import (
	"context"
	"homework/internal/domain"
	"homework/internal/storage"
	"log"

	"github.com/google/uuid"
)

type UserStorage interface {
	Persist(ctx context.Context, userDB storage.UserDB) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Find(ctx context.Context, id uuid.UUID) (domain.UserOutput, error)
	GetList(ctx context.Context, limit, offset int) ([]domain.UserOutput, error)
}

type UserServices struct {
	store UserStorage
}

func NewUserServices(store UserStorage) *UserServices {
	return &UserServices{store: store}
}

func (s *UserServices) Persist(ctx context.Context, userReq CreateUserInput) (string, error) {
	userDB := storage.UserDB{
		Login:        userReq.Login,
		PasswordHash: HashPassword(userReq.Password),
		Name:         userReq.Name,
		Email:        userReq.Email,
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

func (s *UserServices) Find(ctx context.Context, userID string) (domain.UserOutput, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return domain.UserOutput{}, err
	}

	user, err := s.store.Find(ctx, id)
	if err != nil {
		return domain.UserOutput{}, err
	}

	return ToUserResponse(user), nil
}

func (s *UserServices) GetList(ctx context.Context, limit, offset int) ([]handler.UserResponse, error) {
	users, err := s.store.GetList(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return ToUserListResponse(users), nil
}

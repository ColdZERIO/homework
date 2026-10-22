package services

import (
	"context"
	"homework/internal/auth"
	"os"
)

func (s *Service) AddRefreshToken(ctx context.Context, userID, tokenString string) error {
	err := s.storage.AddRefreshToken(ctx, userID, tokenString)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	token, err := s.storage.GetRefreshToken(ctx, userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) CheckRefreshToken(tokenString string) error {
	_, err := auth.ValidateToken(tokenString, []byte(os.Getenv("REFRESH_KEY")))
	if err != nil {
		return err
	}

	return nil
}

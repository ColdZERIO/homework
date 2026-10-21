package services

import "context"

func (s *UserServices) AddRefreshToken(ctx context.Context, userID string, tokenString string)  error {
	err := s.storage.AddRefreshToken(ctx, userID, tokenString)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServices) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	token, err := s.storage.GetRefreshToken(ctx, userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

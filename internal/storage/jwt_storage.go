package storage

import (
	"context"
	"log/slog"
)

func (s *Storage) UpdateRefreshToken(ctx context.Context, userID string, tokenString string) error {
	var token UserTokenModel

	err := s.db.WithContext(ctx).Model(&token).Where("id = ?", userID).Update("refresh_token", tokenString).Error
	if err != nil {
		s.logger.Error(
			"failed to add token",
			slog.String("user_id", userID),
			slog.Any("error", err),
		)

		return err
	}

	return nil
}

func (s *Storage) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	var token UserTokenModel

	err := s.db.WithContext(ctx).First(&token, userID).Error
	if err != nil {
		s.logger.Error(
			"failed to get refresh token",
			slog.String("user_id", userID),
			slog.Any("error", err),
		)

		return "", err
	}

	return token.RefreshToken, nil
}

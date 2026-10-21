package storage

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Storage) AddRefreshToken(ctx context.Context, userID string, tokenString string) error {
	var token = UserTokenModel{
		ID:           uuid.New().String(),
		UserID:       userID,
		RefreshToken: tokenString,
		ExpiredAt:    time.Now().Add(7 * 24 * time.Hour),
		IsActive:     true,
	}

	err := s.db.WithContext(ctx).Transaction(func(db *gorm.DB) error {
		err := db.Model(&token).Where("user_id = ?", userID).Update("is_active = ?", false).Error
		if err != nil {
			return err
		}

		err = s.db.WithContext(ctx).Create(&token).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error(
			"failed to add new refresh token",
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

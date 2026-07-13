package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userID int, login string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"login":   login,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (int, string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("sig method not allowed")
		}

		return secret, nil
	})

	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("invalid user_id")
	}

	login, ok := claims["login"].(string)
	if !ok {
		return 0, "", errors.New("invalid login")
	}

	return int(userID), login, nil
}

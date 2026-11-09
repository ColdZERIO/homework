package auth

import (
	"errors"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v4"
)

type jwtType string

const (
	AccessToken  jwtType = "access"
	RefreshToken jwtType = "refresh"
)

type Claims struct {
	UserID string  `json:"user_id"`
	Role   string  `json:"role"`
	Type   jwtType `json:"type"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("ACCESS_KEY")
	if secret == "" {
		return "", errors.New("ACCESS_KEY is required")
	}
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("REFRESH_KEY")
	if secret == "" {
		return "", errors.New("REFRESH_KEY is required")
	}
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, expectedType jwtType, secret []byte) (*Claims, error) {
	if len(secret) == 0 {
		return nil, errors.New("token secret is required")
	}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrInvalidKey
			}

			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if claims.Type != expectedType {
		return nil, jwt.ErrInvalidKeyType
	}
	if claims.UserID == "" || claims.ExpiresAt == nil {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func ValidateToken(tokenString string, secret []byte) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("invalid signing method")
			}

			return secret, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

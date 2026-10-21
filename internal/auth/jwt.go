package auth

import (
	"errors"
	"os"
	"time"

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
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("ACCESS_KEY")))
}

func GenerateRefreshToken(userID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("REFRESH_KEY")))
}

func ParseToken(tokenString string, expectedType jwtType, secret []byte) (*Claims, error) {
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

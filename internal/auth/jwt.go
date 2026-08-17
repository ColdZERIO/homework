package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
	RoleUser         = "user"
	RoleAdmin        = "admin"
)

type Claims struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role,omitempty"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenManager(secret string, accessTTL, refreshTTL time.Duration) (*TokenManager, error) {
	if secret == "" {
		return nil, errors.New("SECRET_KEY is empty")
	}
	return &TokenManager{secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}, nil
}

func (m *TokenManager) GeneratePair(userID, role string) (TokenPair, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil || parsedID == uuid.Nil {
		return TokenPair{}, errors.New("invalid user ID")
	}

	accessToken, err := m.generate(parsedID.String(), role, TokenTypeAccess, m.accessTTL)
	if err != nil {
		return TokenPair{}, fmt.Errorf("generate access token: %w", err)
	}
	refreshToken, err := m.generate(parsedID.String(), "", TokenTypeRefresh, m.refreshTTL)
	if err != nil {
		return TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}
	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (m *TokenManager) generate(userID, role, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID, Role: role, TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userID, ID: uuid.NewString(),
			IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *TokenManager) ParseAccessToken(tokenString string) (*Claims, error) {
	return m.parse(tokenString, TokenTypeAccess)
}

func (m *TokenManager) ParseRefreshToken(tokenString string) (*Claims, error) {
	return m.parse(tokenString, TokenTypeRefresh)
}

func (m *TokenManager) parse(tokenString, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return m.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	if claims.TokenType != expectedType {
		return nil, errors.New("invalid token type")
	}
	if claims.ExpiresAt == nil {
		return nil, errors.New("token expiration is missing")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil || userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}
	if tokenID, err := uuid.Parse(claims.ID); err != nil || tokenID == uuid.Nil {
		return nil, errors.New("invalid token ID")
	}

	claims.UserID = userID.String()
	return claims, nil
}

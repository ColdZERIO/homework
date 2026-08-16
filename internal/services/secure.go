package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	return string(hash), nil
}

func ValidationPassword(inputPass, hashPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(inputPass))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}

package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func ValidationPassword(inputPass, hashPass string) error {
	checkPass := HashPassword(inputPass)

	if checkPass != hashPass {
		return errors.New("incorrect password")
	}

	return nil
}

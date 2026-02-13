package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) (string, error) {
	if len(p) < 8 {
		return "", errors.New("Password must be at least 8 characters")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

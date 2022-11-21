package security

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashed), nil
}

func VerifyPassword(hashedPass string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
}

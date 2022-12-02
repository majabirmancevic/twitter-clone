package security

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"unicode"
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

func IsValid(s string) bool {
	var (
		hasMinLen = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
		//hasSpecial = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
			//case strings.ContainsAny(s, "<>*()/") == false:
			//	hasSpecial = true
		}
	}

	result := hasMinLen && hasUpper && hasLower && hasNumber
	log.Println("REZULTAT BEZ PROVERE KARAKTERA JE ", result)

	if strings.ContainsAny(s, "<>*()/^#$%&") == true {
		result = false
		log.Println("REZULTAT POSLE JE ", result)
	}
	return result
	//hasMinLen && hasUpper && hasLower && hasNumber && strings.ContainsAny(s, "<>*()/^#$%&") == true
	//&& hasSpecial
}

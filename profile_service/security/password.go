package security

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/mail"
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
		}
	}

	result := hasMinLen && hasUpper && hasLower && hasNumber
	log.Println("REZULTAT BEZ PROVERE KARAKTERA JE ", result)

	if strings.ContainsAny(s, "<>*()/^#$%&;|") == true {
		result = false
		log.Println("REZULTAT POSLE JE ", result)
	}
	return result
}

func Valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidString(s string) bool {

	for _, char := range s {

		if (unicode.IsLetter(char) == true) && (strings.ContainsAny(s, "<>*()/[]") == false) && (strings.Contains(s, "SELECT") == false) && (strings.Contains(s, "FROM") == false) && (strings.Contains(s, "WHERE") == false) {
			return true
		}
	}
	return false
}

func VerifyBusinessInputs(companyName string, email string, webSite string, username string, password string) bool {
	if IsValidString(companyName) && Valid(email) && IsValidString(webSite) && IsValidString(username) && IsValid(password) {
		return true
	}
	return false
}

func VerifyInputs(name string, lastaname string, placeOfLiving string, username string, password string, email string, gender string, age int32) bool {

	if IsValidString(name) && IsValidString(lastaname) && IsValidString(placeOfLiving) && IsValidString(username) && IsValid(password) && Valid(email) && (IsValidString(gender) && (strings.Contains(gender, "M") || strings.Contains(gender, "F"))) && (age >= 13) {
		return true
	}
	return false
}

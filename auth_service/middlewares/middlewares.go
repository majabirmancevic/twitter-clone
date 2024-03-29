package middlewares

import (
	"auth-service/security"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := ExtractToken(r)
		if err != nil {
			security.WriteError(w, http.StatusUnauthorized, security.ErrUnauthorized)
			return
		}
		token, err := ParseToken(tokenString)
		if err != nil {
			log.Println("error on parse token:", err.Error())
			security.WriteError(w, http.StatusUnauthorized, security.ErrUnauthorized)
			return
		}
		if !token.Valid {
			log.Println("invalid token", tokenString)
			security.WriteError(w, http.StatusUnauthorized, security.ErrUnauthorized)
			return
		}

		next(w, r)
	})
}

func ExtractToken(r *http.Request) (string, error) {
	// Authorization => Bearer Token...
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	splitted := strings.Split(header, " ")
	if len(splitted) != 2 {
		log.Println("error on extract token from header:", header)
		return "", errors.New("invalid jwt")
	}
	return splitted[1], nil
}

func parseJwtCallback(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return jwtSecretKey, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, parseJwtCallback)
}

package middlewares

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"

	//"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

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

type TokenPayload struct {
	Username  string
	Role      string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewTokenPayload(tokenString string) (*TokenPayload, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		return nil, errors.New("invalid jwt")
	}
	username, _ := claims["iss"].(string)
	role, _ := claims["aud"].(string)
	createdAt, _ := claims["iat"].(float64)
	expiresAt, _ := claims["exp"].(float64)
	return &TokenPayload{
		Username:  username,
		Role:      role,
		CreatedAt: time.Unix(int64(createdAt), 0),
		ExpiresAt: time.Unix(int64(expiresAt), 0),
	}, nil
}

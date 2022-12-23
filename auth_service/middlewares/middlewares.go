package middlewares

import (
	"auth-service/security"
	"log"
	"net/http"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := security.ExtractToken(r)
		if err != nil {
			security.WriteError(w, http.StatusUnauthorized, security.ErrUnauthorized)
			return
		}
		token, err := security.ParseToken(tokenString)
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
	}
}

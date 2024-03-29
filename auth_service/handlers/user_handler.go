package handlers

import (
	"auth-service/model"
	"auth-service/repository"
	"auth-service/security"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type KeyUser struct{}

type AuthHandler struct {
	logger *log.Logger
	// NoSQL: injecting auth repository
	repo *repository.AuthRepo
}

func NewAuthHandler(l *log.Logger, r *repository.AuthRepo) *AuthHandler {
	return &AuthHandler{l, r}
}

func (p *AuthHandler) SignIn(rw http.ResponseWriter, h *http.Request) {
	log.Println("--------Provera response.Body-------- ")
	if h.Body == nil {
		security.WriteError(rw, http.StatusBadRequest, security.ErrEmptyBody)
		return
	}

	credentials := h.Context().Value(KeyUser{}).(*model.SignInRequest)
	log.Println("--------Provera kredencijala : ", credentials)
	if credentials == nil {
		security.WriteError(rw, http.StatusBadRequest, errors.New("password and username are required "))
		return
	}

	if (security.IsValidString(credentials.Username) && security.IsValid(credentials.Password)) == false {
		security.WriteError(rw, http.StatusBadRequest, errors.New("credentials are not valid "))
		return
	}

	log.Println("--------Provera da li korisnik postoji-------- ")
	user, err := p.repo.GetByUsername(credentials.Username)
	log.Println("--------GetByUsername ERROR : ", err)
	if err != nil {
		log.Println("This user does not exist")
		log.Println(err.Error())

		security.WriteError(rw, http.StatusBadRequest, errors.New("sign in failed - This user does not exist"))
		return
	}

	if !user.Verified {
		security.WriteError(rw, http.StatusForbidden, errors.New("sign in failed - This user does not verified"))
		return
	}

	log.Println("--------Provera passworda-------- ")
	err = security.VerifyPassword(user.Password, credentials.Password)
	log.Println("--------VerifyPassword ERROR : ", err)
	if err != nil {
		log.Println("Passwords does not matches !!!", err.Error())
		security.WriteError(rw, http.StatusForbidden, errors.New("sign in failed - Passwords does not matches"))
		return
	}

	log.Println("--------Kreiranje tokena-------- ")
	token, err := security.NewToken(user.Username, user.Role)
	log.Println("--------NewToken ERROR : ", err)
	if err != nil {
		log.Println("Token cannot be created", err.Error())
		security.WriteError(rw, http.StatusForbidden, errors.New("sign in failed - Token cannot be created"))
		return
	}

	response := model.SignInResponseRegular{
		Token:    token,
		Username: user.Username,
		//RegularProfile: model.NewUserResponse(user),
	}

	security.WriteAsJson(rw, http.StatusOK, response)

}

func (p *AuthHandler) MiddlewareLoginDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		req := &model.SignInRequest{}
		err := req.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyUser{}, req)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (p *AuthHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		p.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, h)

	})
}

func (p *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {

	var tokenStr string
	p.logger.Println("STRING ZA PROVERU : ", tokenStr)

	err := json.NewDecoder(r.Body).Decode(&tokenStr)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//tokenString, err := security.ExtractToken(r)
	//if err != nil {
	//	security.WriteError(w, http.StatusUnauthorized, security.ErrUnauthorized)
	//	return
	//}
	token, err := security.ParseToken(tokenStr)
	if err != nil {
		log.Println("error on parse token:", err.Error())
		security.WriteError(w, http.StatusUnauthorized, security.ErrUnauthorized)
		return
	}
	if !token.Valid {
		log.Println("invalid token", tokenStr)
		security.WriteError(w, http.StatusUnauthorized, security.ErrUnauthorized)
		return
	}

	tokenPayload, err := security.NewTokenPayload(tokenStr)
	if err != nil {
		log.Println("cant generate token payload :", err.Error())
		security.WriteError(w, http.StatusUnauthorized, security.ErrUnauthorized)
		return
	}

	if tokenPayload.Role != "regular" || tokenPayload.Role != "business" {
		log.Println("permission error ", tokenPayload.Role)
		security.WriteError(w, http.StatusUnauthorized, security.ErrUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

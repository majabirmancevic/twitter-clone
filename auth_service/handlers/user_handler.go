package handlers

import (
	"auth-service/model"
	"auth-service/repository"
	"auth-service/security"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"net/mail"
	"strings"
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

func (p *AuthHandler) SignUp(rw http.ResponseWriter, h *http.Request) {
	user := h.Context().Value(KeyUser{}).(*model.RegularProfile)
	user.Email = strings.ToLower(user.Email)
	user.ID = primitive.NewObjectID()
	hashedPassword, _ := security.EncryptPassword(user.Password)
	user.Password = hashedPassword

	found, _ := p.repo.GetByUsername(user.Username)

	if found != nil && found.Username == user.Username {
		http.Error(rw, "User already exist with this username !", http.StatusBadRequest)
		p.logger.Fatal("This username is already used !")
	}

	if found == nil || (found != nil && found.Username != user.Username) {

		if valid(user.Email) == false {
			http.Error(rw, "Email is not valid !", http.StatusBadRequest)
			p.logger.Println("NEISPRAVNA EMAIL ADRESA !")
		}
		p.logger.Println("------- SLANJE U BAZU")
		p.repo.Insert(user)
		rw.WriteHeader(http.StatusCreated)
		p.logger.Println("------- USPESNO KREIRAN KORISNIK")

	}

	//rw.WriteHeader(http.StatusBadRequest)
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

	log.Println("--------Provera da li korisnik postoji-------- ")
	user, err := p.repo.GetByUsername(credentials.Username)
	log.Println("--------GetByUsername ERROR : ", err)
	if err != nil {
		log.Println("This user does not exist")
		log.Println(err.Error())
		errors.New("sign in failed - This user does not exist")
	}

	log.Println("--------Provera passworda-------- ")
	err = security.VerifyPassword(user.Password, credentials.Password)
	log.Println("--------VerifyPassword ERROR : ", err)
	if err != nil {
		log.Println("Passwords does not matches !!!", err.Error())
		errors.New("sign in failed - Passwords does not matches")
	}

	log.Println("--------Kreiranje tokena-------- ")
	token, err := security.NewToken(user.ID.Hex())
	log.Println("--------NewToken ERROR : ", err)
	if err != nil {
		log.Println("Token cannot be created", err.Error())
		errors.New("sign in failed - Token cannot be created")
	}

	response := model.SignInResponseRegular{
		Token:          token,
		RegularProfile: model.NewUserResponse(user),
	}

	security.WriteAsJson(rw, http.StatusOK, response)
}

func (p *AuthHandler) GetAllRegularUsers(rw http.ResponseWriter, h *http.Request) {
	users, err := p.repo.GetAll()
	if err != nil {
		p.logger.Print("Database exception: ", err)
	}

	if users == nil {
		return
	}

	err = users.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *AuthHandler) MiddlewareUserDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		user := &model.RegularProfile{}
		err := user.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyUser{}, user)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
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

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

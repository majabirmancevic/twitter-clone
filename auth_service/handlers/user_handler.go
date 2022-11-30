package handlers

import (
	"auth-service/model"
	"auth-service/repository"
	"auth-service/security"
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"
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
	user.Role = "regular"
	user.Verified = false

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
		// Generate Verification Code
		code := randstr.String(20)
		verificationCode := security.Encode(code)
		p.logger.Println("------- VERIFIKACIONI KOD KREIRAN ", verificationCode)

		user.VerificationCode = verificationCode

		p.logger.Println("------- SLANJE U BAZU")
		err := p.repo.Insert(user)
		if err != nil {
			http.Error(rw, "Neuspesno dodavanje korisnika !", http.StatusBadRequest)
			p.logger.Println(" ----- Error ", err)
		}

		var firstName = user.Name

		if strings.Contains(firstName, " ") {
			firstName = strings.Split(firstName, " ")[1]
		}

		// 👇 Send Email
		emailData := security.EmailData{
			URL:       "Your account verification code is " + code,
			FirstName: firstName,
			Subject:   "Account verification",
		}

		p.logger.Println("------- slanje mejla ", emailData)
		security.SendEmail(user, &emailData)

		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Origin", "*")

		rw.WriteHeader(http.StatusCreated)
		p.logger.Println("------- USPESNO KREIRAN KORISNIK")
		return

	}

	//rw.WriteHeader(http.StatusBadRequest)
}

// 	VERIFY EMAIL

func (p *AuthHandler) VerifyEmail(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	code := vars["code"]
	verificationCode := security.Encode(code)

	log.Println("CODE ", code)
	log.Println("Verification code ", verificationCode)
	//var updatedUser model.RegularProfile
	updatedUser, err := p.repo.GetByVerificationCode(verificationCode)

	log.Println("UPDATED USER ", updatedUser)

	if err != nil {
		log.Println("Invalid verification code or user doesn't exists")
		log.Println(err.Error())
		response := errors.New("Invalid verification code or user doesn't exists")
		security.WriteAsJson(rw, http.StatusBadRequest, response)
		return
	}

	if updatedUser.Verified {
		log.Println("User already verified")
		log.Println(err.Error())
		response := errors.New("User already verified")
		security.WriteAsJson(rw, http.StatusConflict, response)
		return
	}

	updatedUser.VerificationCode = ""
	updatedUser.Verified = true
	p.repo.Update(updatedUser.ID, updatedUser)

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(http.StatusOK)

	//ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
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
		//rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		//rw.Header().Set("Access-Control-Allow-Origin", "*")

		//rw.Header().Set("Content-Type", "text/html; charset=utf-8")

		next.ServeHTTP(rw, h)
	})
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

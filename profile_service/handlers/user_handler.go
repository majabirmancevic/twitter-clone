package handlers

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"profile_service/model"
	"profile_service/repository"
	"profile_service/security"
	"strings"
)

type KeyUser struct{}
type BusinessUser struct{}

type ProfileHandler struct {
	logger *log.Logger
	// NoSQL: injecting auth repository
	repo *repository.ProfileRepo
}

func NewProfileHandler(l *log.Logger, r *repository.ProfileRepo) *ProfileHandler {
	return &ProfileHandler{l, r}
}

func (p *ProfileHandler) SignUp(rw http.ResponseWriter, h *http.Request) {
	user := h.Context().Value(KeyUser{}).(*model.RegularProfile)

	log.Println("STIGAO JE ZAHTEV ", h.Method)

	user.Gender = strings.ToUpper(user.Gender)

	if security.VerifyInputs(user.Name, user.Lastname, user.PlaceOfLiving, user.Username, user.Password, user.Email, user.Gender, user.Age) == false {

		p.logger.Println(user.Name, user.Lastname, user.PlaceOfLiving, user.Username, user.Password, user.Email, user.Gender, user.Age)
		p.logger.Println(security.VerifyInputs(user.Name, user.Lastname, user.PlaceOfLiving, user.Username, user.Password, user.Email, user.Gender, user.Age))

		http.Error(rw, "your data input isn't valid", http.StatusInternalServerError)
		return
	}

	if security.CheckBlacklistedPassword(user.Password) {
		p.logger.Fatal("This password is unsafe  !")
		security.WriteAsJson(rw, http.StatusBadRequest, "This password is unsafe !")
		return
	}

	hashedPassword, _ := security.EncryptPassword(user.Password)
	user.Password = hashedPassword
	user.Email = strings.ToLower(user.Email)
	user.ID = primitive.NewObjectID()
	user.Role = "regular"
	user.Verified = false

	found, _ := p.repo.GetByUsername(user.Username)

	if found != nil && found.Username == user.Username {
		p.logger.Fatal("This username is already used !")
		security.WriteAsJson(rw, http.StatusBadRequest, "User already exist with this username !")
		return
	}

	code := randstr.String(20)
	verificationCode := security.Encode(code)
	p.logger.Println("------- VERIFIKACIONI KOD KREIRAN ", verificationCode)

	user.VerificationCode = verificationCode

	p.logger.Println("------- SLANJE U BAZU")
	error := p.repo.Insert(user)
	if error != nil {
		p.logger.Println(" ----- Error ", error)
		security.WriteAsJson(rw, http.StatusBadRequest, "Neuspesno dodavanje korisnika !")
		return
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

	if security.SendEmail(user, &emailData) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Success"))
		p.logger.Println("Uspesno slanje maila", emailData)
	} else {
		security.WriteAsJson(rw, http.StatusInternalServerError, errors.New("something went wrong"))
		p.logger.Println("Neuspesno slanje maila", emailData)
	}

}

func (p *ProfileHandler) SignUpBusiness(rw http.ResponseWriter, h *http.Request) {

	user := h.Context().Value(BusinessUser{}).(*model.BusinessProfile)
	p.logger.Println("BIZNIS USER ", user.CompanyName, user.WebSite, user.Email)
	user.Email = strings.ToLower(user.Email)

	if security.VerifyBusinessInputs(user.CompanyName, user.Email, user.WebSite, user.Username, user.Password) == false {

		security.WriteAsJson(rw, http.StatusBadRequest, errors.New("your data input isn't valid"))
		return
	}

	if security.CheckBlacklistedPassword(user.Password) {
		p.logger.Fatal("This password is unsafe  !")
		security.WriteAsJson(rw, http.StatusBadRequest, "This password is unsafe !")
		return
	}

	hashedPassword, _ := security.EncryptPassword(user.Password)
	user.Password = hashedPassword
	user.ID = primitive.NewObjectID()
	user.Role = "business"
	user.Verified = false

	foundRegular, _ := p.repo.GetByUsername(user.Username)
	foundBusiness, _ := p.repo.GetBusinessByUsername(user.Username)
	if (foundRegular != nil && foundRegular.Username == user.Username) || (foundBusiness != nil && foundBusiness.Username == user.Username) {
		p.logger.Fatal("This username is already used !")
		security.WriteAsJson(rw, http.StatusBadRequest, "User already exist with this username !")
		return
	}

	code := randstr.String(20)
	verificationCode := security.Encode(code)
	p.logger.Println("------- VERIFIKACIONI KOD KREIRAN ", verificationCode)

	user.VerificationCode = verificationCode

	p.logger.Println("------- SLANJE U BAZU")
	error := p.repo.InsertBusinessUser(user)
	if error != nil {
		p.logger.Println(" ----- Error ", error)
		security.WriteAsJson(rw, http.StatusBadRequest, "Neuspesno dodavanje korisnika !")
		return
	}

	var firstName = user.CompanyName

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
	security.SendEmailBusiness(user, &emailData)

	rw.WriteHeader(http.StatusCreated)
}

// -------------------------------------------------------------------------------------------------
func (p *ProfileHandler) VerifyEmail(rw http.ResponseWriter, h *http.Request) {
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

	rw.WriteHeader(http.StatusOK)

}

func (p *ProfileHandler) VerifyBusinessEmail(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	code := vars["code"]
	verificationCode := security.Encode(code)

	log.Println("CODE ", code)
	log.Println("Verification code ", verificationCode)
	//var updatedUser model.RegularProfile
	updatedUser, err := p.repo.GetBusinessByVerificationCode(verificationCode)

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
	p.repo.UpdateBusiness(updatedUser.ID, updatedUser)

	rw.WriteHeader(http.StatusOK)

}

//-------------------------------------------------------------------------------------------

func (p *ProfileHandler) GetRegularUser(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	username := vars["username"]

	user, err := p.repo.GetByUsername(username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	response := model.NewUserResponse(user)
	security.WriteAsJson(rw, http.StatusOK, response)

}

func (p *ProfileHandler) GetBusinessUser(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	username := vars["username"]

	user, err := p.repo.GetBusinessByUsername(username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//response := model.NewUserResponse(user)
	security.WriteAsJson(rw, http.StatusOK, user)

}

//------------------------------------------------------------------------------------------

func (p *ProfileHandler) GetAllRegularUsers(rw http.ResponseWriter, h *http.Request) {
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

func (p *ProfileHandler) GetAllBusinessUsers(rw http.ResponseWriter, h *http.Request) {
	users, err := p.repo.GetAllBusiness()
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

//-----------------------------------------------------------------------------------------------

func (p *ProfileHandler) DeleteAll(rw http.ResponseWriter, h *http.Request) {

	p.repo.DeleteAll()
	rw.WriteHeader(http.StatusNoContent)
}

// ------------------------------------------------------------------------------------------
func (p *ProfileHandler) MiddlewareUserDeserialization(next http.Handler) http.Handler {
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

func (p *ProfileHandler) MiddlewareBusinessUserDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		user := &model.BusinessProfile{}
		err := user.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), BusinessUser{}, user)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (p *ProfileHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		p.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")
		//rw.Header().Add("Access-Control-Allow-Headers", "Content-Type,Origin,Content-Type, Accept, Authorization")
		//rw.Header().Add("Access-Control-Allow-Origin", "*")

		//if h.Method == "OPTIONS" {
		//	rw.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH")
		//	rw.WriteHeader(http.StatusOK)
		//}
		//rw.Header().Set("Content-Type", "text/html; charset=utf-8")

		next.ServeHTTP(rw, h)

	})
}

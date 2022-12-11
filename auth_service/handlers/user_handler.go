package handlers

import (
	"auth-service/model"
	"auth-service/repository"
	"auth-service/security"
	"context"
	"crypto/tls"
	"errors"
	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
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

	log.Println("STIGAO JE ZAHTEV ", h.Method)

	user.Gender = strings.ToUpper(user.Gender)

	if security.VerifyInputs(user.Name, user.Lastname, user.PlaceOfLiving, user.Username, user.Password, user.Email, user.Gender, user.Age) == false {

		p.logger.Println(user.Name, user.Lastname, user.PlaceOfLiving, user.Username, user.Password, user.Email, user.Gender, user.Age)
		p.logger.Println(security.VerifyInputs(user.Name, user.Lastname, user.PlaceOfLiving, user.Username, user.Password, user.Email, user.Gender, user.Age))

		http.Error(rw, "your data input isn't valid", http.StatusInternalServerError)
		return
	}

	//p.logger.Println("PROVERA U BLEK LISTI ", security.CheckBlacklistedPassword(user.Password))
	//passwords, err := security.LoadPasswords()
	//log.Print(len(passwords))
	//p.logger.Println("PRVIH 5 IZ BLEK LISTE ", passwords[0], passwords[1], passwords[2], passwords[3], passwords[4], passwords[5])
	//if err != nil {
	//	http.Error(rw, "Can't read file", http.StatusInternalServerError)
	//}

	if security.CheckBlacklistedPassword(user.Password) {
		p.logger.Fatal("This password is unsafe  !")
		http.Error(rw, "This password is unsafe !", http.StatusBadRequest)
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
		http.Error(rw, "User already exist with this username !", http.StatusBadRequest)
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
		http.Error(rw, "Neuspesno dodavanje korisnika !", http.StatusBadRequest)
		return
	}

	var firstName = user.Name
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	//p.logger.Println("------- slanje mejla ")
	//security.SendMailSMTP(code, firstName)

	// 👇 Send Email
	emailData := security.EmailData{
		URL:       "Your account verification code is " + code,
		FirstName: firstName,
		Subject:   "Account verification",
	}

	p.logger.Println("------- slanje mejla ", emailData)

	if security.SendEmail(user, &emailData) {
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte("Success"))
		p.logger.Println("Uspesno slanje maila", emailData)
	} else {
		http.Error(rw, "Neuspesno dodavanje korisnika !", http.StatusInternalServerError)
		p.logger.Println("------- poslat 500 response ", emailData)
	}

}

func (p *AuthHandler) SignUpBusiness(rw http.ResponseWriter, h *http.Request) {

	user := h.Context().Value(KeyUser{}).(*model.BusinessProfile)
	user.Email = strings.ToLower(user.Email)

	if security.VerifyBusinessInputs(user.CompanyName, user.Email, user.WebSite, user.Username, user.Password) == false {

		http.Error(rw, "your data input isn't valid", http.StatusBadRequest)
		return
	}

	//p.logger.Println("PROVERA U BLEK LISTI ", security.CheckBlacklistedPassword(user.Password))
	//passwords, err := security.LoadPasswords()
	//log.Print(len(passwords))
	//p.logger.Println("PRVIH 5 IZ BLEK LISTE ", passwords[0], passwords[1], passwords[2], passwords[3], passwords[4], passwords[5])
	//if err != nil {
	//	http.Error(rw, "Can't read file", http.StatusInternalServerError)
	//	return
	//}

	if security.CheckBlacklistedPassword(user.Password) {
		p.logger.Fatal("This password is unsafe  !")
		http.Error(rw, "This password is unsafe !", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := security.EncryptPassword(user.Password)
	user.Password = hashedPassword
	user.ID = primitive.NewObjectID()
	user.Role = "business"
	user.Verified = false

	found, _ := p.repo.GetByUsername(user.Username)

	if found != nil && found.Username == user.Username {
		p.logger.Fatal("This username is already used !")
		http.Error(rw, "User already exist with this username !", http.StatusBadRequest)
		return
	}

	code := randstr.String(20)
	verificationCode := security.Encode(code)
	p.logger.Println("------- VERIFIKACIONI KOD KREIRAN ", verificationCode)

	user.VerificationCode = verificationCode

	p.logger.Println("------- SLANJE U BAZU")
	error := p.repo.InsertBusiness(user)
	if error != nil {
		p.logger.Println(" ----- Error ", error)
		http.Error(rw, "Neuspesno dodavanje korisnika !", http.StatusBadRequest)
		return
	}

	var firstName = user.CompanyName
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	//p.logger.Println("------- slanje mejla ")
	//if security.SendMailSMTP(code, firstName) {
	//
	//}

	//// 👇 Send Email
	//emailData := security.EmailData{
	//	URL:       "Your account verification code is " + code,
	//	FirstName: firstName,
	//	Subject:   "Account verification",
	//}
	//
	//p.logger.Println("------- slanje mejla ", emailData)
	//security.SendEmail(user, &emailData)

	rw.WriteHeader(http.StatusCreated)
}

func (p *AuthHandler) SendingMailTest(w http.ResponseWriter, h *http.Request) {
	//from := "ibsit2022@gmail.com"
	//password := "xeuloaiprwagrouh"
	//
	//toEmailAddress := "ibsit2022@gmail.com"
	//to := []string{toEmailAddress}
	//
	//host := "smtp.gmail.com"
	//port := "587"
	//address := host + ":" + port

	code := randstr.String(20)
	emailData := security.EmailData{
		URL: "Your account verification code is " + code,
		//FirstName: firstName,
		Subject: "Account verification",
	}

	//subject := "To: ibsit2022@gmail.com\r\n" + "Account verification \r\n"
	//body := "Your account verification code is " + code
	//message := []byte(subject + body)
	//
	//auth := smtp.PlainAuth("", from, password, host)
	//log.Println("AUTH => ", auth)
	//
	//err := smtp.SendMail(address, auth, from, to, message)
	//log.Println("ERROR => ", err)
	//if err != nil {
	//	panic(err)
	//	log.Println(" ---- Greska -> ", err.Error())
	//	http.Error(w, "Neuspesno slanje mejla!", http.StatusBadRequest)
	//	return
	//}

	from := "ibsit2022@gmail.com"
	smtpPass := "xeuloaiprwagrouh"
	smtpUser := "ibsit2022@gmail.com"
	to := []string{"ibsit2022@gmail.com"}
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to[0])
	m.SetHeader("Subject", emailData.Subject)
	m.SetBody("text/html", emailData.URL)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.TLSConfig = &tls.Config{}
	dial, err := d.Dial()
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	dial.Send(from, to, m)

	log.Println("HOST-> ", d.Host, " PORT-> ", d.Port, " LOCALNAME-> ", d.LocalName, " SSL-> ", d.SSL, " TLSCONFIG-> ", d.TLSConfig)

	// Send Email
	//if err := d.DialAndSend(m); err != nil {
	//	log.Fatal("Could not send email: ", err)
	//}

	w.WriteHeader(http.StatusOK)
	dial.Close()
	log.Println(" ---- Uspesno poslat mejl  ")
}

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
	token, err := security.NewToken(user.Username)
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

func (p *AuthHandler) MiddlewareBusinessUserDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		user := &model.BusinessProfile{}
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

func (p *AuthHandler) DeleteAll(rw http.ResponseWriter, h *http.Request) {

	p.repo.DeleteAll()
	rw.WriteHeader(http.StatusNoContent)
}

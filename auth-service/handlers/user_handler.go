package handlers

import (
	"auth-service/model"
	"auth-service/security"
	"auth-service/service"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type UserHandler struct {
	service *service.UserService
	cli     *mongo.Client
	logger  *log.Logger
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) Init(r *mux.Router) {
	r.Use(handler.MiddlewareContentTypeSet)
	loginRouter := r.Methods(http.MethodPost).Subrouter()
	loginRouter.HandleFunc("/login", handler.UserLogin)
	loginRouter.Use(handler.MiddlewareContentTypeSet)

	signupRouter := r.Methods(http.MethodPost).Subrouter()
	signupRouter.HandleFunc("/signup", handler.UserSignup)
	signupRouter.Use(handler.MiddlewareContentTypeSet)
	//r.HandleFunc("/signup", handler.UserSignup).Methods("POST")
	http.Handle("/", r)
}

func (handler *UserHandler) UserSignup(response http.ResponseWriter, request *http.Request) {
	username := request.URL.Query().Get("username")

	findUser, error := handler.service.GetByUsername(username)
	if error != nil {
		handler.logger.Print("Database exception: ", error)
	}
	if findUser.Username == username {
		handler.logger.Print("User already exists with this username: ", error)
	}

	var user model.RegularProfile
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	user.ID = primitive.NewObjectID()
	result, _ := handler.service.SignUp(&user)
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(result)

}

func (handler *UserHandler) UserLogin(response http.ResponseWriter, request *http.Request) {
	var req model.SignInRequest
	var dbUser model.RegularProfile

	json.NewDecoder(request.Body).Decode(&req)
	collection := handler.cli.Database("twitter").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&dbUser)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	userPass := []byte(req.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		log.Println(passErr)
		response.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}
	jwtToken, err := security.GenerateJWT()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	response.Write([]byte(`{"token":"` + jwtToken + `"}`))

}

func (handler *UserHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		handler.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}

package main

import (
	"context"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"profile_service/handlers"
	"profile_service/repository"
	"time"
)

func main() {

	//Reading from environment, if not set we will default it to 8080.
	//This allows flexibility in different environments (for eg. when running multiple docker api's and want to override the default port)
	port := os.Getenv("PROFILE_PORT")
	if len(port) == 0 {
		port = "8002"
	}

	// Initialize context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[auth-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[auth-store] ", log.LstdFlags)

	// NoSQL: Initialize Product Repository store
	store, err := repository.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Disconnect(timeoutContext)

	// NoSQL: Checking if the connection was established
	store.Ping()

	userHandler := handlers.NewProfileHandler(logger, store)

	router := mux.NewRouter()
	//router.Use(userHandler.MiddlewareContentTypeSet)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", userHandler.SignUp)
	postRouter.Use(userHandler.MiddlewareUserDeserialization)

	postBusinessRouter := router.Methods(http.MethodPost).Subrouter()
	postBusinessRouter.HandleFunc("/business", userHandler.SignUpBusiness)
	postBusinessRouter.Use(userHandler.MiddlewareBusinessUserDeserialization)

	verifyRouter := router.Methods(http.MethodGet).Subrouter()
	verifyRouter.HandleFunc("/verifyEmail/{code}", userHandler.VerifyEmail)

	getUserRouter := router.Methods(http.MethodGet).Subrouter()
	getUserRouter.HandleFunc("/user/{username}", userHandler.GetRegularUser)

	// ZA PROVERU PRISTUPA RUTA NA OSNOVU TOKENA
	//middlewares.Authenticate(userHandler.SignIn)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedHeaders([]string{"Origin, Content-Type, X-Auth-Token"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH"}))

	//Initialize the server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)

	go func() {
		err := server.ListenAndServeTLS("auth_service/certificates/self-ssl.crt", "auth_service/certificates/self-ssl.key")
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")

}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"tweet_service/data"
	"tweet_service/handlers"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	//Reading from environment, if not set we will default it to 8080.
	//This allows flexibility in different environments (for eg. when running multiple docker api's and want to override the default port)
	port := os.Getenv("TWEET_PORT")
	if len(port) == 0 {
		port = "8003"
	}

	// Initialize context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[tweet-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[tweet-store] ", log.LstdFlags)

	// NoSQL: Initialize Product Repository store
	store, err := data.New(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.CloseSession()

	store.CreateTables()

	//Initialize the handler and inject said logger
	tweetsHandler := handlers.NewTweetsHandler(logger, store)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()
	router.Use(tweetsHandler.MiddlewareContentTypeSet)

	getTweetIds := router.Methods(http.MethodGet).Subrouter()
	getTweetIds.HandleFunc("/tweets/users/{tweetId}", tweetsHandler.GetAllRegularUserIds)

	getTweetByRegUser := router.Methods(http.MethodGet).Subrouter()
	getTweetByRegUser.HandleFunc("/tweets/{username}", tweetsHandler.GetTweetsByRegUser)

	postTweetForRegUser := router.Methods(http.MethodPost).Subrouter()
	postTweetForRegUser.HandleFunc("/tweets", tweetsHandler.CraeteTweetForRegUser)
	postTweetForRegUser.Use(tweetsHandler.MiddlewareTweetForRegUserDeserialization)

	likeTweetForRegUser := router.Methods(http.MethodPost).Subrouter()
	likeTweetForRegUser.HandleFunc("/like/{id}", tweetsHandler.LikeTweet)
	likeTweetForRegUser.Use(tweetsHandler.MiddlewareLikeDeserialization)

	getLikesByTweet := router.Methods(http.MethodGet).Subrouter()
	getLikesByTweet.HandleFunc("/likes/{id}", tweetsHandler.GetLikesByTweet)

	deleteLikesByTweet := router.Methods(http.MethodDelete).Subrouter()
	deleteLikesByTweet.HandleFunc("/dislike/{tweetId}/{username}", tweetsHandler.DislikeTweet)

	getCountLikesByTweet := router.Methods(http.MethodGet).Subrouter()
	getCountLikesByTweet.HandleFunc("/likes/count/{id}", tweetsHandler.GetCountByLikes)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	//Initialize the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
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

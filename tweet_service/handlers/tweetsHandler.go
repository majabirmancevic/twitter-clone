package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tweet_service/data"
)

type KeyTweet struct{}

type TweetsHandler struct {
	logger *log.Logger
	// NoSQL: injecting student repository
	repo *data.TweetRepo
}

// Injecting the logger makes this code much more testable.
func NewTweetsHandler(l *log.Logger, r *data.TweetRepo) *TweetsHandler {
	return &TweetsHandler{l, r}
}

func (s *TweetsHandler) GetAllRegularUserIds(rw http.ResponseWriter, h *http.Request) {
	tweetIds, err := s.repo.GetDistinctIds("regular_username", "tweet_by_regular_user")
	if err != nil {
		s.logger.Print("Database exception: ", err)
	}

	if tweetIds == nil {
		return
	}

	s.logger.Println(tweetIds)

	e := json.NewEncoder(rw)
	err = e.Encode(tweetIds)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		s.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (s *TweetsHandler) LikeTweet(rw http.ResponseWriter, h *http.Request){


}

func (s *TweetsHandler) GetTweetsByRegUser(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	username := vars["username"]

	tweetsByRegUser, err := s.repo.GetTweetsByUser(username)
	if err != nil {
		s.logger.Print("Database exception: ", err)
	}

	if tweetsByRegUser == nil {
		return
	}

	err = tweetsByRegUser.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		s.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (s *TweetsHandler) CraeteTweetForRegUser(rw http.ResponseWriter, h *http.Request) {
	userTweet := h.Context().Value(KeyTweet{}).(*data.TweetByRegularUser)
	err := s.repo.InsertTweetByRegUser(userTweet)
	if err != nil {
		s.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (s *TweetsHandler) MiddlewareTweetForRegUserDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		tweet := &data.TweetByRegularUser{}
		err := tweet.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			s.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyTweet{}, tweet)
		h = h.WithContext(ctx)
		next.ServeHTTP(rw, h)
	})
}

func (s *TweetsHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		s.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}

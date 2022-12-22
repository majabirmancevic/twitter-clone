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
	vars := mux.Vars(h)
	tweetId := vars["tweetId"]
	tweetIds, err := s.repo.GetDistinctIds("username", "likes_by_user", tweetId)
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

func (s *TweetsHandler) LikeTweet(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	log.Println(" >>>>>> Tweet id : ", id)
	payload := h.Context().Value(KeyTweet{}).(*data.Like)

	likes, err := s.repo.GetLikesByTweet(id)
	if err != nil {
		s.logger.Print("Error : ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, like := range likes {
		currentUsername := payload.Username
		if currentUsername == like.Username {
			s.logger.Print("You have already liked this tweet ")
			http.Error(rw, "You have already liked this tweet", http.StatusBadRequest)
			//rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	error := s.repo.InsertLikeByRegUser(payload.Username, id)
	if error != nil {
		s.logger.Print("Database exception (INSERT LIKE): ", error)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
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

func (s *TweetsHandler) GetLikesByTweet(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	likesByTweet, err := s.repo.GetLikesByTweet(id)
	if err != nil {
		s.logger.Print("Database exception: ", err)
		s.logger.Print("Lajkovi : ", len(likesByTweet))
	}

	if likesByTweet == nil {
		return
	}

	err = likesByTweet.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		s.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (s *TweetsHandler) GetCountByLikes(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	count, _ := s.repo.GetLikesByTweet(id)
	likesByTweet := len(count)

	if likesByTweet == -1 {
		return
	}

	err := data.ToJSON(rw, likesByTweet)
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

func (s *TweetsHandler) DislikeTweet(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	username := vars["username"]
	tweetId := vars["tweetId"]

	err := s.repo.DeleteLikeByUser(tweetId, username)
	if err != nil {
		s.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
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

func (s *TweetsHandler) MiddlewareLikeDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		like := &data.Like{}
		err := like.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			s.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyTweet{}, like)
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

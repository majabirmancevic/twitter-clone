package service

import (
	"auth-service/model"
	"auth-service/security"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type UserService struct {
	store model.UserStore
}

func NewUserService(store model.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) Get(id string) (*model.RegularProfile, error) {
	return s.store.GetById(id)
}

func (s *UserService) GetByUsername(username string) (*model.RegularProfile, error) {
	return s.store.GetByUsername(username)
}

func (s *UserService) SignUp(req *model.RegularProfile) (*model.RegularProfile, error) {
	req.Password = security.EncryptPassword(req.Password)
	req.Name = strings.TrimSpace(req.Name)
	req.Lastname = strings.TrimSpace(req.Lastname)
	req.Username = strings.TrimSpace(req.Username)
	req.Tweets = []model.Tweet{}
	found, err := s.store.GetByUsername(req.Username)
	if err == mongo.ErrNoDocuments {
		user := new(model.RegularProfile)
		err := s.store.Insert(user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	if found == nil {
		return nil, err
	}

	return nil, errors.New("email already exists")
}

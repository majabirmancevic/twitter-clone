package service

import (
	"auth-service/repository"
)

type UserService struct {
	store repository.UserRepository
}

func NewUserService(store repository.UserRepository) *UserService {
	return &UserService{
		store: store,
	}
}

//func (service *UserService) GetById(id primitive.ObjectID) (*model.RegularProfile, error) {
//	return service.store.GetById(id)
//}

//func (service *UserService) Save(ctx context.Context, user *model.RegularProfile)

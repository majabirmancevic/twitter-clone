package model

type UserStore interface {
	GetById(id string) (*RegularProfile, error)
	GetByUsername(username string) (*RegularProfile, error)
	Insert(order *RegularProfile) error
}

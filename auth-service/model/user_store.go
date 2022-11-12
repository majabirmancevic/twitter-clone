package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserStore interface {
	Get(id primitive.ObjectID) (*RegularProfile, error)
	Insert(order *RegularProfile) error
}

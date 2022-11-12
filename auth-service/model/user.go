package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegularProfile struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string
	Lastname      string
	Gender        string
	Age           int32
	PlaceOfLiving string
	Username      string //mora biti jedinstveno
	Password      string
	IsPrivate     bool
	Tweets        []Tweet
}

type BusinessProfile struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CompanyName string
	Email       string
	WebSite     string
	Username    string //mora biti jedinstveno
	Password    string
	Tweets      []Tweet
}

type Tweet struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

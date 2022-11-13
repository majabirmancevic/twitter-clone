package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegularProfile struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name,omitempty" json:"name"`
	Lastname      string             `bson:"lastname,omitempty" json:"lastname"`
	Gender        string             `bson:"gender,omitempty" json:"gender"`
	Age           int32              `bson:"age,omitempty" json:"age"`
	PlaceOfLiving string             `bson:"placeOfLiving,omitempty" json:"placeOfLiving"`
	//mora biti jedinstveno
	Username  string  `bson:"username,omitempty" json:"username"`
	Password  string  `bson:"password,omitempty" json:"password"`
	IsPrivate bool    `bson:"isPrivate" json:"isPrivate"`
	Tweets    []Tweet `bson:"tweets" json:"tweets"`
}

type BusinessProfile struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CompanyName string             `bson:"companyName,omitempty" json:"companyName"`
	Email       string             `bson:"email,omitempty" json:"email"`
	WebSite     string             `bson:"webSite,omitempty" json:"webSite"`
	//mora biti jedinstveno
	Username string  `bson:"username,omitempty" json:"username"`
	Password string  `bson:"password,omitempty" json:"password"`
	Tweets   []Tweet `bson:"tweets,omitempty" json:"tweets"`
}

type Tweet struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

type SignUpRegularRequest struct {
	Name          string `bson:"name,omitempty" json:"name"`
	Lastname      string `bson:"lastname,omitempty" json:"lastname"`
	Gender        string `bson:"gender,omitempty" json:"gender"`
	Age           uint16 `bson:"age,omitempty" json:"age"`
	PlaceOfLiving string `bson:"placeOfLiving,omitempty" json:"placeOfLiving"`
	Username      string `bson:"username,omitempty" json:"username"`
	Password      string `bson:"password,omitempty" json:"password"`
}

type SignInRequest struct {
	Username string `bson:"username,omitempty" json:"username"`
	Password string `bson:"password,omitempty" json:"password"`
}

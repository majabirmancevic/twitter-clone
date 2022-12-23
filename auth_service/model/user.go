package model

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type RegularProfiles []*RegularProfile

// Model for create and read user from db
type RegularProfile struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username" validate:"required"`
	Password string             `bson:"password" json:"password" validate:"required,min=8"`
	Verified bool               `bson:"verified" json:"verified" `
	Role     string             `json:"role" bson:"role"`
}

type SignInRequest struct {
	Username string `bson:"username" json:"username" binding:"required"`
	Password string `bson:"password" json:"password" binding:"required"`
}

type SignInResponseRegular struct {
	Token    string `bson:"token" json:"token" `
	Username string `bson:"username" json:"username" `
}

//---------------------------------------------------------------
//
//func (p *RegularProfiles) ToJSON(w io.Writer) error {
//	e := json.NewEncoder(w)
//	return e.Encode(p)
//}

func (p *RegularProfile) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *RegularProfile) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
func (p *SignInRequest) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

package twitter_clone

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	ID       primitive.ObjectID
	Name     string
	Surname  string
	Sex      string
	Age      string
	Adress   Address
	Email    string
	Username string
}

type Address struct {
	Street  string
	City    string
	Country string
}

type Profiles []*Profile

func (p *Profile) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Profile) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

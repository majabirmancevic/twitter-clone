package data

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"io"
)

type TweetByRegularUser struct {
	RegularUsername string
	Description     string
	Id              gocql.UUID
}

type TweetsByRegularUser []*TweetByRegularUser

func (o *TweetsByRegularUser) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *TweetByRegularUser) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *TweetByRegularUser) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

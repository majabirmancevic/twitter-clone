package data

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"io"
)

type TweetByRegularUser struct {
	LikeCounter     uint64     `json:"like_counter"`
	RegularUsername string     `json:"regular_username"`
	Description     string     `json:"description"`
	Id              gocql.UUID `json:"id"`
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

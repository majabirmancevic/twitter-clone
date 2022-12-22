package data

import (
	"encoding/json"
	"io"
)

type TweetByRegularUser struct {
	RegularUsername string `json:"regular_username"`
	Description     string `json:"description"`
	//Id              gocql.UUID `json:"id"`
	Id string `json:"id"`
}

type Like struct {
	Username string `json:"username"`
	TweetId  string `json:"tweetId"`
	Id       string `json:"likeId"`
}

type TweetsByRegularUser []*TweetByRegularUser
type Likes []*Like

func (o *Like) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Like) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *TweetsByRegularUser) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *TweetByRegularUser) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Likes) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func ToJSON(w io.Writer, broj int) error {
	e := json.NewEncoder(w)
	return e.Encode(broj)
}

func (o *TweetByRegularUser) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

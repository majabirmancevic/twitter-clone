package security

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrEmptyBody    = errors.New("body can't be empty")
	ErrUnauthorized = errors.New("Unauthorized")
)

type JError struct {
	Error string `json:"error"`
}

func WriteAsJson(w http.ResponseWriter, statusCode int, data interface{}) {
	//w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	e := "error"
	if err != nil {
		e = err.Error()
	}
	WriteAsJson(w, statusCode, JError{e})
}

func Encode(s string) string {
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return string(data)
}

func Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

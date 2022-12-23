package security

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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

func CheckBlacklistedPassword(password string) bool {
	passwords, err := LoadPasswords()
	if err != nil {
		return true
	}

	for _, passwordCheck := range passwords {
		if password == passwordCheck {
			return true
		}
	}

	return false
}

func LoadPasswords() ([]string, error) {

	//path, _ := filepath.Abs("./auth_service/security/password_blacklist.txt")
	//log.Println("ABS ", path)

	file, err := os.Open("./auth_service/security/password_blacklist.txt")
	if err != nil {
		log.Fatal("ERROR ", err)
		return []string{}, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var passwords []string

	for scanner.Scan() {
		passwords = append(passwords, scanner.Text())
		//log.Println("FOR PETLJA ", len(passwords))
	}

	if err = file.Close(); err != nil {
		fmt.Printf("Could not close the file due to this %s error \n", err)
	}

	return passwords, nil
}

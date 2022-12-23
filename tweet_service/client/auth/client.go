package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Client struct {
	address string
}

func NewClient() *Client {
	return &Client{
		address: fmt.Sprintf("https://localhost:8001/verifyToken"),
	}
}

func (client *Client) VerifyToken(token string) error {

	reqBytes, err := json.Marshal(token)
	log.Println("---> TOKEN MARSHAL ", reqBytes, "ERROR ", err)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(reqBytes)
	requestURL := client.address
	log.Println("---> requestURL ", requestURL)
	httpReq, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)

	log.Println("---> httpReq ", httpReq)
	if err != nil {
		log.Println("---> httpReq error : ", err)
		return errors.New("error updating inventory")
	}

	res, err := http.DefaultClient.Do(httpReq)

	if err != nil || res.StatusCode != http.StatusOK {
		log.Println("---> Error conn :", err)
		log.Println("---> Status code : ", res.StatusCode)
		return errors.New("error updating inventory")
	}
	return nil
}

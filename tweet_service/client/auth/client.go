package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Client struct {
	address string
}

//func NewClient() *Client {
//	return &Client{
//		address: fmt.Sprintf("https://localhost:8000/auth_service/verifyToken"),
//	}
//}

func VerifyToken(client *http.Client, method string, token string) error {

	endpoint := "https://localhost:8000/auth_service/verifyToken"
	reqBytes, err := json.Marshal(token)

	bodyReader := bytes.NewReader(reqBytes)
	//requestURL := client.address
	//og.Println("---> requestURL ", requestURL)
	httpReq, err := http.NewRequest(method, endpoint, bodyReader)
	httpReq.Header.Set("Content-Type", "application/json")

	log.Println("---> httpBody : ", httpReq.Body, "---> httpHeader : ", httpReq.Header, "---> httpURL : ", httpReq.URL)
	if err != nil {
		log.Println("---> httpReq error : ", err)
		return errors.New(err.Error())
	}

	res, err := client.Do(httpReq)
	//log.Printf("RESPONSE CODE %v & STATUS %v", res.StatusCode, res.Status)

	if err != nil || res.StatusCode != http.StatusOK {
		log.Println("---> Error conn :", err)
		//log.Println("---> Status code : ", res.StatusCode)
		return errors.New(err.Error())
	}

	//// Close the connection to reuse it
	//defer res.Body.Close()
	//
	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	log.Fatalf("Couldn't parse response body. %+v", err)
	//	return nil
	//}

	return nil
}

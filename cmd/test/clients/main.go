package main

import (
	"log"
	"net/http"
)

import "bytes"

import "io/ioutil"

func LoginRequest() (string, error) {
	var buf = []byte(`{"Username":"legacy", "Password":"blah"}`)
	req, err := http.NewRequest("POST", "http://localhost:3000/clients/login", bytes.NewBuffer(buf))
	if err != nil {
		log.Print("error with request: ", err)
		return "err", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("error with sending request", err)
		return "err", err
	}
	defer resp.Body.Close()
	log.Print("token generated")
	body, _ := ioutil.ReadAll(resp.Body)
	tokenString := bytes.NewBuffer(body).String()
	return tokenString, nil
}

func CharacterCreateRequest(token string) (string, error) {

	var characterCreateJSON = []byte(`{"Token":"` + token + `", "Name" : "legacy33"}`)
	req, err := http.NewRequest("POST", "http://localhost:3000/characters/new", bytes.NewBuffer(characterCreateJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Error with request 2: ", err)
		return "err", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Print("create character response: ", string(body))
	return string(body), nil
}

func DisconnectRequest(token string) (string, error) {

	buf := []byte(token)
	req, err := http.NewRequest("POST", "http://localhost:3000/clients/disconnect", bytes.NewBuffer(buf))
	if err != nil {
		log.Print("error with request: ", err)
		return "err", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("error with sending request", err)
		return "err", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Print("disonnect response: ", string(body))
	return string(body), nil
}

func main() {
	//time.Sleep(time.Minute * 2)

	token, err := LoginRequest()
	if err != nil {
		log.Print("error sending login request", err)
	}

	_, err = CharacterCreateRequest(token)
	if err != nil {
		log.Print("error sending create character request", err)
	}

	_, err = DisconnectRequest(token)
	if err != nil {
		log.Print("error sending disconnect request", err)
	}
}

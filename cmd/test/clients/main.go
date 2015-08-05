package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"thorium-go/requests"
	"time"
)

import "bytes"

import "io/ioutil"

func LoginRequest(username string, password string) (request.LoginResponse, error) {
	var loginReq request.Authentication
	var martiniResponse request.LoginResponse
	loginReq.Username = username
	loginReq.Password = password
	jsonBytes, err := json.Marshal(&loginReq)
	if err != nil {
		return martiniResponse, err
	}

	req, err := http.NewRequest("POST", "http://localhost:6960/clients/login", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Print("error with request: ", err)
		return martiniResponse, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("error with sending request", err)
		return martiniResponse, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&martiniResponse)
	if err != nil {
		log.Print("bad json request", resp.Body)
		return martiniResponse, err
	}

	return martiniResponse, nil
}

func CharacterSelectRequest(token string, id int) (string, error) {

	var selectReq request.SelectCharacter
	selectReq.AccountToken = token
	selectReq.ID = id
	jsonBytes, err := json.Marshal(&selectReq)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:6960/characters/%d/select", id), bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Error with request 2: ", err)
		return "err", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func CharacterCreateRequest(token string, name string) (string, error) {

	var charCreateReq request.CreateCharacter
	charCreateReq.AccountToken = token
	charCreateReq.Name = name
	jsonBytes, err := json.Marshal(&charCreateReq)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "http://localhost:6960/characters/new", bytes.NewBuffer(jsonBytes))
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
	req, err := http.NewRequest("POST", "http://localhost:6960/clients/disconnect", bytes.NewBuffer(buf))
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

	resp, err := LoginRequest("legacy", "blah")
	if err != nil {
		log.Print("error sending login request", err)
	}
	log.Print("LoginResponse Token: ", resp.UserToken)
	log.Print("LoginResponse Character ID's: ", resp.CharacterIDs)

	var charSession string
	var selected bool = false
	for i := 0; i < len(resp.CharacterIDs); i++ {
		log.Print("character %d id = %d", i, resp.CharacterIDs[i])

		if resp.CharacterIDs[i] != 0 {
			charSession, err = CharacterSelectRequest(resp.UserToken, resp.CharacterIDs[i])
			if err != nil {
				log.Print(err)
				continue
			}

			selected = true
			break
		}
	}

	if !selected {
		rand.Seed(int64(time.Now().Second()))

		charSession, err = CharacterCreateRequest(resp.UserToken, fmt.Sprintf("legacy%d", rand.Intn(100000)))
		if err != nil {
			log.Print(err)
		}
	}

	log.Print("character session:\n", charSession)

	_, err = DisconnectRequest(resp.UserToken)
	if err != nil {
		log.Print("error sending disconnect request", err)
	}
}

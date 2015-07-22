package main

import (
	"net/http"
	"thorium-go/requests"
)
import "bytes"
import "fmt"
import "io/ioutil"

func main() {
	var buf = []byte(`{"Username":"legacy", "Password":"blah"}`)
	req, err := http.NewRequest("POST", "http://localhost:3000/clients/login", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println("error with request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error with sending request", err)
	}
	defer resp.Body.Close()
	fmt.Println("token generated")
	var charCreateReq request.CreateCharacter
	body, _ := ioutil.ReadAll(resp.Body)
	charCreateReq.Token = bytes.NewBuffer(body).String()

	//time.Sleep(time.Minute * 2)

	var tokenB = []byte(`{"Token":"` + charCreateReq.Token + `", "Name" : "legacy"}`)
	req2, err := http.NewRequest("POST", "http://localhost:3000/characters/new", bytes.NewBuffer(tokenB))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := client.Do(req2)
	if err != nil {
		fmt.Println("Error with request 2: ", err)
	}
	defer resp2.Body.Close()
	body2, _ := ioutil.ReadAll(resp2.Body)
	fmt.Println("response Body:", string(body2))

	buf = []byte(charCreateReq.Token)
	req, err = http.NewRequest("POST", "http://localhost:3000/clients/disconnect", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println("error with request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("error with sending request", err)
	}
	defer resp.Body.Close()

}

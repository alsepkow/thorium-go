package main

import "net/http"
import "bytes"
import "fmt"
import "io/ioutil"
import "hussain/thorium-go/requests"

func main() {
	var buf = []byte(`{"Username":"legacy", "Password":"test"}`)
	req, err := http.NewRequest("POST", "http://localhost:3000/client/login", bytes.NewBuffer(buf))
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
	var SomeRequest request.Test
	body, _ := ioutil.ReadAll(resp.Body)
	SomeRequest.Token = bytes.NewBuffer(body).String()

	//time.Sleep(time.Minute * 2)

	var tokenB = []byte(`{"Token":"` + SomeRequest.Token + `"}`)
	req2, err := http.NewRequest("POST", "http://localhost:3000/client/afterlogin", bytes.NewBuffer(tokenB))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := client.Do(req2)
	if err != nil {
		fmt.Println("Error with request 2: ", err)
	}
	defer resp2.Body.Close()
	body2, _ := ioutil.ReadAll(resp2.Body)
	fmt.Println("response Body:", string(body2))

}

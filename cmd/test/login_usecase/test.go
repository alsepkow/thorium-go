package main

import (
	"log"
	"thorium-go/client"
)

func main() {
	//time.Sleep(time.Minute * 2)

	userToken, err := client.LoginRequest("legacy", "blah")
	if err != nil {
		log.Print("error sending login request", err)
	}
	/*
		//	chars := make([10]int)
		//	_, err = ViewCharacters(&chars)
		if err != nil {
			log.Print(err)
		}
		// foreach character data print it
		// here

		// use this when done above
		//_, err = CharacterSelectRequest(token, chars[0])
		var charSession string
		charSession, err = client.CharacterSelectRequest(token, 6)
	*/
	var charToken string
	charToken, err = client.CharacterSelectRequest(userToken, 6)
	//	charToken, err = client.CharacterCreateRequest(userToken, "legacy33")
	if err != nil {
		log.Print("error sending create character request", err)
	}

	log.Print("character session:\n", charToken)

	_, err = client.DisconnectRequest(userToken)
	if err != nil {
		log.Print("error sending disconnect request", err)
	}
}

package main

import (
	"log"
	"thorium-go/database"
)

func main() {
	log.Print("Creating Character")
	token, _ := thordb.LoginAccount("legacy", "blah")
	log.Print("logged in")
	thordb.Disconnect(token)
	log.Print("disconnected")
}

package main

import (
	"encoding/json"
	"log"
	"thorium-go/database"
)

func main() {

	testCharacterJson()

}

func testCharacterJson() {

	var character thordb.Character
	b, err := json.Marshal(&character)
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(string(b))
}

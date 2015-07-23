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
	character := thordb.NewCharacter()
	b, err := json.Marshal(character)
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(string(b))
	b, err = json.Marshal(character.GameData)
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(string(b))

	character.ID = 2
	character.UserID = 5
	// public only for tsting
	_, err = thordb.StoreCharacterSnapshot(character)
	if err != nil {
		log.Print(err)
	} else {
		log.Print("success")
	}
}

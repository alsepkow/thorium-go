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
	charSession := thordb.NewCharacterSession()
	b, err := json.Marshal(charSession.CharacterData)
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(string(b))

	// using legacy33 test char
	charSession.ID = 2
	charSession.UserID = 5
	// this needs to require a machine token later
	_, err = thordb.StoreCharacterSnapshot(charSession)
	if err != nil {
		log.Print(err)
	} else {
		log.Print("success")
	}
}

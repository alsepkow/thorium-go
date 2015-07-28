package thordb

import "log"

func TestMachineRequest() {
	_, err := kvstore.Ping().Result()
	if err != nil {
		log.Print(err)
	}

	log.Print("Made it!")
}

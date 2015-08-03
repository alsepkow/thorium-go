package thordb

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const machineSessionKey string = "machines/%d"
const hkeyMachineToken string = "machineToken"

func RegisterMachine(remoteAddress string, servicePort int) (int, string, error) {

	var machineId int
	err := db.QueryRow("INSERT INTO machines (remote_address, service_listen_port) VALUES ($1, $2) RETURNING machine_id", remoteAddress, servicePort).Scan(&machineId)
	if err != nil {
		return 0, "", err
	}

	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims["machineId"] = machineId
	token.Claims["iat"] = time.Now()
	var token_str string
	token_str, err = token.SignedString(signKey)
	if err != nil {
		return 0, "", err
	}

	_, err = db.Exec("INSERT INTO machines_metadata  VALUES ($1, $2, $3, $4, $5, $6)", machineId, token_str, time.Now(), 0.0, 0.0, 0.0)
	if err != nil {
		return 0, "", err
	}

	var ok bool
	ok, err = kvstore.HSet(fmt.Sprintf(machineSessionKey, machineId), hkeyMachineToken, token_str).Result()
	if err != nil {
		return 0, "", err
	}
	if !ok {
		return 0, "", errors.New("thordb: unable to set machine token in redis")
	}
	kvstore.Expire(fmt.Sprintf(machineSessionKey, machineId), time.Second*120)

	return machineId, token_str, nil
}

func UnregisterMachine(machineToken string) (bool, error) {

	machineId, err := validateMachineToken(machineToken)
	if err != nil {
		return false, err
	}

	res, err := db.Exec("DELETE FROM machines WHERE machine_id = $1", machineId)
	if err != nil {
		log.Print("couldn't delete machine from postgres")
	} else {
		var rows int64
		rows, err = res.RowsAffected()
		if err != nil {
			log.Print("couldnt read rows affected")
		} else if rows == 0 {
			log.Print("couldn't delete machine from postgres - does not exist")
		}
	}

	var count int64
	count, err = kvstore.Del(fmt.Sprintf(machineSessionKey, machineId)).Result()
	if err != nil {
		log.Print("couldn't delete machine from redis cache")
		log.Print(err)
	}

	if count == 0 {
		log.Print("couldn't delete machine from redis cache")
	}
	return true, nil
}

func UpdateMachineStatus(machineToken string, usageCpu float64, usageNetwork float64, usagePlayerCapacity float64) error {

	machineId, err := validateMachineToken(machineToken)
	if err != nil {
		return err
	}

	// ToDo: use this later to check that 1 row was updated
	//var res sql.Result
	_, err = db.Exec("UPDATE machines_metadata SET last_heartbeat = $1, cpu_usage_pct = $2, network_usage_pct = $3, player_occupancy_pct = $4 WHERE machine_id = $5",
		time.Now(), usageCpu, usageNetwork, usagePlayerCapacity, machineId)
	if err != nil {
		return err
	}

	kvstore.Expire(fmt.Sprintf(machineSessionKey, machineId), time.Second*120)

	return nil
}

func TestMachineRequest() {
	_, err := kvstore.Ping().Result()
	if err != nil {
		log.Print(err)
	}

	log.Print("Made it!")
}

func validateMachineToken(token_str string) (int, error) {
	token, err := jwt.Parse(token_str, func(t *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return 0, err
	}

	var idFloat64 float64
	idFloat64, ok := token.Claims["machineId"].(float64)
	id := int(idFloat64)
	if !ok {
		return 0, errors.New("thordb: invalid session")
	}

	var savedToken string
	savedToken, err = kvstore.HGet(fmt.Sprintf(machineSessionKey, id), hkeyMachineToken).Result()
	if err != nil {
		return 0, err
	}
	if token_str == savedToken {
		return id, nil
	} else {
		return 0, errors.New("thordb: invalid session")
	}
}

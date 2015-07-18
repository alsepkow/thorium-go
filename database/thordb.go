package thordb

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

var db *sql.DB
var kvstore *redis.Client

func init() {
	var err error
	db, err = sql.Open("postgres", "user=thoriumnet password=thoriumtest dbname=thoriumnet host=localhost")
	if err != nil {
		panic(err)
	}

	kvstore = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err = kvstore.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func RegisterNewGame(mapName string, maxPlayers int) (int, error) {
	var gameId int
	err := db.QueryRow("INSERT INTO games (map_name, max_players) VALUES ( $1, $2 ) RETURNING game_id", mapName, maxPlayers).Scan(&gameId)
	return gameId, err
}

func RegisterActiveGame(gameId int, machineId int, port int) (bool, error) {
	res, err := db.Exec("SELECT * FROM public.games WHERE game_id = $1", gameId)
	if err != nil {
		return false, err
	}

	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		return false, err
	}

	exists := rows > 0
	if exists != true {
		log.Print("gameId ", strconv.Itoa(gameId), " does not exist")
		return false, err
	}

	res, err = db.Exec("INSERT INTO active_games (game_id, machine_id, port) VALUES ( $1, $2, $3 )", gameId, machineId, port)
	if err != nil {
		return false, err
	}

	rows, err = res.RowsAffected()
	if err != nil {
		return false, err
	}

	exists = rows > 0
	return exists, err
}

func RegisterMachine(machineIp string, port int) (int, error) {
	var machineId int
	err := db.QueryRow("INSERT INTO game_machines (ip_address, port) VALUES ( $1, $2 ) RETURNING machine_id", machineIp, port).Scan(&machineId)
	return machineId, err
}

func UnregisterMachine(machineId int) (bool, error) {

	success := false

	res, err := db.Exec("DELETE FROM game_machines WHERE machine_id = $1", machineId)
	if err != nil {
		return success, err
	}

	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		return success, err
	}

	success = rows > 0
	return success, err

}

func RegisterAccount(username string, password string) (int, error) {
	var foundname string
	err := db.QueryRow("SELECT username FROM account_data WHERE username LIKE $1;", username).Scan(&foundname)
	switch {
	case err == sql.ErrNoRows:
		log.Print("Username available")
	case err != nil:
		log.Print(err)
		return 0, err
	default:
		log.Print("Username is already in use")
		return 0, errors.New("thordb: already in use")
	}

	saltSize := 16
	var alg string = "sha1"
	//allocates 16+sha1.Size bytes to the bufer
	//creates slice with length saltSize and capacity of saltSize+sha1.Size
	buf := make([]byte, saltSize, saltSize+sha1.Size)

	//fill buf with random data (linux is /dev/urandom)
	_, e := io.ReadFull(rand.Reader, buf)

	if e != nil {
		fmt.Println("filling buf with random data failed")
		return 0, e
	}

	dirtySalt := sha1.New()
	dirtySalt.Write(buf)
	dirtySalt.Write([]byte(password))
	salt := dirtySalt.Sum(buf)

	combination := string(salt) + string(password)
	passwordHash := sha1.New()
	io.WriteString(passwordHash, combination)
	log.Printf("Password Hash : %x \n", passwordHash.Sum(nil))
	log.Printf("BSalt: %x \n", salt)
	var uid int
	var timenow = time.Now()
	err = db.QueryRow("INSERT INTO account_data (username, password, salt, algorithm, createdon, lastlogin) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id", username, passwordHash.Sum(nil), salt, alg, timenow, timenow).Scan(&uid)
	if err != nil {
		fmt.Println("error inserting account data: ", err)
		return 0, err
	}

	return uid, err
}

func LoginAccount(username string, password string) (string, error) {
	var token string

	// get the account info from db
	var hashedPassword []byte
	var salt []byte
	var uid int
	err := db.QueryRow("SELECT password, salt, user_id FROM account_data WHERE username LIKE $1", username).Scan(&hashedPassword, &salt, &uid)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("thordb: user does not exist %s", username)
		return token, errors.New("thordb: does not exist")
	case err != nil:
		log.Print(err)
		return token, err
	}

	combination := string(salt) + string(password)
	passwordHash := sha1.New()
	io.WriteString(passwordHash, combination)
	match := bytes.Equal(passwordHash.Sum(nil), hashedPassword)
	if match {
		// register a new session in redis
		token = "OK"
		return token, nil
	} else {
		return token, errors.New("thordb: invalid password")
	}
}

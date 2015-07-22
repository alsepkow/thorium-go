package thordb

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

const privKeyPath string = "keys/app.rsa"
const pubKeyPath string = "keys/app.rsa.pub"

// redis keys
const accountSessionKey string = "sessions/account/%d"
const characterSessionKey string = "sessions/character/%d"
const gameSessionKey string = "games/%d"

var db *sql.DB
var kvstore *redis.Client
var signKey *rsa.PrivateKey
var verifyKey *rsa.PublicKey

func init() {
	// check rsa
	var signBytes []byte
	var verifyBytes []byte
	var err error
	log.Print("opening app.rsa keys")
	signBytes, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Print(err)
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Print(err)
	}
	verifyBytes, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Print(err)
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Print(err)
	}

	log.Print("testing postgres connection")
	// check postgres
	db, err = sql.Open("postgres", "user=thoriumnet password=thoriumtest dbname=thoriumnet host=localhost")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("testing redis connection")
	// check redis
	kvstore = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err = kvstore.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("thordb initialization complete")
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
	var token_str string

	// get the account info from db
	var hashedPassword []byte
	var salt []byte
	var uid int
	err := db.QueryRow("SELECT password, salt, user_id FROM account_data WHERE username LIKE $1", username).Scan(&hashedPassword, &salt, &uid)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("thordb: user does not exist %s", username)
		return token_str, errors.New("thordb: does not exist")
	case err != nil:
		log.Print(err)
		return token_str, err
	}

	combination := string(salt) + string(password)
	passwordHash := sha1.New()
	io.WriteString(passwordHash, combination)
	match := bytes.Equal(passwordHash.Sum(nil), hashedPassword)
	if !match {
		return token_str, errors.New("thordb: invalid password")
	}

	// create the jwt token
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims["uid"] = uid
	token.Claims["iat"] = time.Now()
	token_str, err = token.SignedString(signKey)

	if err != nil {
		return token_str, err
	}

	// I'm going to insert the raw token data into redis here, but is that security proof?
	// In future we could maybe use a field in the encrypted claim as the session key? im not sure if that works or not though

	// first check if a session already exists, if so reject as "already logged on" unless the time is substantially old (> 5min)
	var alreadyLoggedIn bool = true

	_, err = kvstore.Get(fmt.Sprintf(accountSessionKey, uid)).Result()
	if err != nil {
		switch err.Error() {
		case "redis: nil":
			alreadyLoggedIn = false
		default:
			return token_str, err
		}
	}

	if alreadyLoggedIn {
		return token_str, errors.New("thordb: already logged in")
	}

	// set the session in redis and give it a 2 minute expiry
	// the client needs to ping once every 2 minutes to refresh the expiry
	kvstore.Set(fmt.Sprintf(accountSessionKey, uid), token_str, 0)
	kvstore.Expire(fmt.Sprintf(accountSessionKey, uid), time.Second*120)
	return token_str, nil
}

func Disconnect(accountSessionToken string) error {

	token, err := jwt.Parse(accountSessionToken, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		return err
	}

	log.Printf("thordb: disconnect w/ token: \n%+v\n", token)

	var uidFloat64 float64
	uidFloat64, ok := token.Claims["uid"].(float64)
	uid := int(uidFloat64)
	if !ok {
		log.Print("couldnt convert uid")
		log.Print("actual type", reflect.TypeOf(token.Claims["uid"]))
		return errors.New("thordb: invalid session")
	}

	// ToDo: update account + character in postgres before deleting from redis

	log.Printf("thordb: disconnect uid = %d", uid)
	var count int64
	count, err = kvstore.Del(fmt.Sprintf(accountSessionKey, uid)).Result()
	if err != nil {
		return err
	}

	if count == 0 {
		log.Print("couldnt find session")
		return errors.New("thordb: invalid session")
	}

	log.Print("client disconnected %d", uid)
	return nil
}

// helper funcs
func storeAccount(session *AccountSession) {
	// store in postgres
}

func storeCharacter(character_session *CharacterSession) {

}

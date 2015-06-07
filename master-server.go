package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

import "database/sql"
import _ "github.com/lib/pq"
import "github.com/go-martini/martini"

type NewGameRequest struct {
	Map        string
	MaxPlayers int
}

func main() {
	fmt.Println("hello world")

	m := martini.Classic()
	m.Post("/games/new_request", handleGameRequest)
	m.Post("/games/:id/register_server", handleRegisterServer)
	//m.Get("/games/:id", handleGetGameInfo)
	//m.Post("/games/:id/heartbeat_server")
	m.Post("/machines/register_new", handleRegisterMachine)
	m.Post("/machines/:id/unregister", handleUnregisterMachine)
	m.Run()
}

func handleRegisterMachine(req *http.Request) (int, string) {

	db, err := connectToDB()
	if err != nil {
		fmt.Println("[ThoriumNET] unable to connect to DB")
		fmt.Println(err)
		return 500, "Internal Server Error"

	}

	machineIp := req.RemoteAddr

	var machineId string
	err = db.QueryRow("INSERT INTO game_machines (ip_endpoint) VALUES ($1) RETURNING machine_id", machineIp).Scan(&machineId)
	if err != nil {
		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	fmt.Println("machine registered, ip=", machineIp)
	return 200, machineId
}

func handleUnregisterMachine(params martini.Params) (int, string) {

	machineId := params["id"]

	db, err := connectToDB()
	if err != nil {
		fmt.Println("[ThoriumNET] unable to connect to DB")
		fmt.Println(err)
		return 500, "Internal Server Error"

	}

	_, err = db.Exec("DELETE FROM game_machines WHERE machine_id = $1", machineId)
	if err != nil {
		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	fmt.Println("machine unregistered, id=", machineId)
	return 200, "OK"
}

func handleGameRequest(req *http.Request) (int, string) {
	fmt.Println("[ThoriumNET] master-server.handleGameRequest")

	decoder := json.NewDecoder(req.Body)
	var request NewGameRequest
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	if request.Map == "" {
		fmt.Println("No Map Name Given")
		return 500, "No Map Name Given"
	}

	db, err := connectToDB()
	if err != nil {
		fmt.Println("[ThoriumNET] unable to connect to DB")
		fmt.Println(err)
		return 500, "Internal Server Error"

	}

	// create a new game with max players = 12 as default
	var gameId int
	err = db.QueryRow("INSERT INTO games (ip_endpoint, map_name, max_players) VALUES ( NULL, $1, $2 ) RETURNING game_id", request.Map, request.MaxPlayers).Scan(&gameId)
	if err != nil {
		fmt.Println("[ThoriumNET] unable to insert new game record")
		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	fmt.Println("[ThoriumNET] new game, id=", strconv.Itoa(gameId))
	return 200, "OK"
}

func handleRegisterServer(params martini.Params) (int, string) {

	db, err := connectToDB()
	if err != nil {
		fmt.Println("[ThoriumNET] unable to connect to DB")
		fmt.Println(err)
		return 500, "Internal Server Error"
	}

	gameId := params["id"]

	fmt.Println("[ThoriumNET] master-server.handleRegisterServer ID=", gameId)

	// query to find out if the game id is validi
	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		return 500, "Internal Server Error"
	}

	// if game exists, update ip with remote ip
	var res sql.Result
	res, err = tx.Exec("SELECT * FROM public.games WHERE game_id = $1", gameId)
	rows, err2 := res.RowsAffected()
	if err != nil || err2 != nil || rows == 0 {
		tx.Rollback()
		fmt.Println("[ThoriumNET] game not found, ID=", gameId)
		fmt.Println(err)
		return 500, "Internal Server Error"
	} else {
		var res sql.Result
		res, err = tx.Exec("UPDATE games SET ip_endpoint = '127.0.1.7:39182' WHERE game_id = $1", gameId)
		rowsAffected, err2 := res.RowsAffected()
		if err != nil || err2 != nil || rowsAffected == 0 {
			tx.Rollback()
			fmt.Println("[ThoriumNET] failed to update game ip endpoint, ID=", gameId)
			fmt.Println(err)
			return 500, "Internal Server Error"

		}

		tx.Commit()
		fmt.Println("Found game ", gameId)
		return 200, "OK"
	}
}

func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=thoriumnet password=thoriumtest dbname=thoriumnet host=localhost")
	if err != nil {
		fmt.Println("database conn err: ", err)
		return nil, err
	}

	return db, nil
}

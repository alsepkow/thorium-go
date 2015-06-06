package main

import (
	"fmt"
	"strconv"
)

import "database/sql"
import _ "github.com/lib/pq"
import "github.com/go-martini/martini"

func main() {
	fmt.Println("hello world")

	m := martini.Classic()
	m.Post("/games/new_request", handleGameRequest)
	m.Post("/games/:id/register_server", handleRegisterServer)
	m.Post("/games/:id/heartbeat_server")
	m.Run()
}

func handleGameRequest() (int, string) {
	fmt.Println("[ThoriumNET] master-server.handleGameRequest")

	db, err := connectToDB()
	if err != nil {
		return 500, "Internal Server Error"
		fmt.Println("[ThoriumNET] unable to connect to DB")
	}

	// create a new game
	var gameId int
	err = db.QueryRow("INSERT INTO games (ip_endpoint) VALUES (NULL) RETURNING game_id").Scan(&gameId)
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
		return 500, "Internal Server Error"
	}

	gameId := params["id"]

	fmt.Println("[ThoriumNET] master-server.handleRegisterServer ID=", gameId)

	var pgId int
	var pgEndpoint string

	err = db.QueryRow("SELECT * FROM public.games WHERE game_id = $1").Scan(&pgId, &pgEndpoint)
	if err != nil {
		return 500, "Internal Server Error"
	} else {
		fmt.Println("Found game ", pgId)
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

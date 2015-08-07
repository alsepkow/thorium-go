package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"
)

func main() {
	log.Print("running a fake bolt server")

	var game_id int
	var service_port int
	var listen_port int

	var map_name string
	var game_mode string

	flag.IntVar(&game_id, "id", 0, "identifies this game within the cluster")
	flag.IntVar(&listen_port, "p", 0, "game server listen port")
	flag.IntVar(&service_port, "s", 0, "machine local service port")
	flag.StringVar(&map_name, "m", "", "game map: tutorial, openworld, sandbox")
	flag.StringVar(&game_mode, "g", "", "game mode: tutorial, openworld, sandbox")

	flag.Parse()

	if game_id == 0 || listen_port == 0 || service_port == 0 {
		log.Fatal("bad arguments")
	}

	endpoint := fmt.Sprintf("http://localhost:%d/games/%d/register_server", service_port, game_id)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal("Die - failed to register")
	}

	m := martini.Classic()
	m.Get("/status", handleStatusRequest)
	m.RunOnAddr(fmt.Sprintf(":%d", listen_port))
}

func handleStatusRequest(httpReq *http.Request) (int, string) {
	return 200, "OK"
}

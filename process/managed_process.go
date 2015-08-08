//package main
package process

import (
	"fmt"
	"log"
	"strconv"
)
import "os"
import "os/exec"

type ManagedProcess struct {
	ApplicationName string
	Process         *os.Process
	GameId          int
	ListenPort      int
	GameMode        string
	MapName         string
	// add more here if needed
}

var process_list map[int]ManagedProcess = make(map[int]ManagedProcess)

func NewGameServer(game_id int, listen_port int, service_port int, game_mode string, map_name string) (*ManagedProcess, error) {
	log.Print("starting new game server")

	cmd := exec.Command("bolt-server",
		"-id", strconv.Itoa(game_id),
		"-p", strconv.Itoa(listen_port),
		"-s", strconv.Itoa(service_port),
		"-m", map_name,
		"-g", game_mode)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	var process ManagedProcess
	process.ApplicationName = "bolt-server.go"
	process.Process = cmd.Process
	process.GameId = game_id
	process.ListenPort = listen_port
	process.GameMode = game_mode
	process.MapName = map_name

	return &process, nil
}

func execute_cmd(command string) (*exec.Cmd, error) {
	cmd := exec.Command(command)
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("User Commmand: ", command)
	}
	return cmd, err
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"thorium-go/requests"
	"time"
)

import _ "github.com/lib/pq"
import "github.com/go-martini/martini"

var registerData request.MachineRegisterResponse

func main() {
	fmt.Println("hello world")

	timeNow := time.Now()
	rand.Seed(int64(timeNow.Second()))
	port := rand.Intn(50000)
	port = port + 10000

	fmt.Println(strconv.Itoa(port), "\n")

	reqData := &request.RegisterMachine{Port: port}
	jsonBytes, err := json.Marshal(reqData)
	if err != nil {
		log.Fatal(err)
	}
	request, err := http.NewRequest("POST", "http://localhost:6960/machines/register", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		log.Print("Error registering with master")
		os.Exit(1)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	json.Unmarshal([]byte(body), &registerData)
	fmt.Println("Registered As Machine#", registerData.MachineId)

	m := martini.Classic()

	m.Get("/ping", handlePingRequest)
	//m.Post("/games/new_request", handleGameRequest)
	//m.Post("/games/:id/register_server", handleRegisterServer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		fmt.Println(<-c)
		shutdown()
		os.Exit(1)
	}()
	defer shutdown()

	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				sendHeartbeat()
			}
		}
	}()

	thisIp := fmt.Sprint("localhost:", strconv.Itoa(port))
	m.RunOnAddr(thisIp)

}

func sendHeartbeat() {
	log.Print("...")
}

func handlePingRequest() (int, string) {
	return 200, "OK"
}

func shutdown() {

	var reqData request.UnregisterMachine
	reqData.MachineToken = registerData.MachineToken
	jsonBytes, err := json.Marshal(&reqData)
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest("POST", fmt.Sprintf("http://localhost:6960/machines/%d/disconnect", registerData.MachineId), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Print("failed to disconnect properly")
		return
	}
}

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

import _ "github.com/lib/pq"
import "github.com/go-martini/martini"

func main() {
	fmt.Println("hello world")

	time := time.Now()
	rand.Seed(int64(time.Second()))
	port := rand.Intn(50000)
	port = port + 10000

	fmt.Println(strconv.Itoa(port), "\n")

	var jsonStr = fmt.Sprint(`{"Port":`, port, `}`)
	request, err := http.NewRequest("POST", "http://localhost:3000/machines/register_new", bytes.NewBuffer([]byte(jsonStr)))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	machineId := string(body)
	fmt.Println("Registered As Machine#", machineId)

	m := martini.Classic()

	m.Get("/ping", handlePingRequest)
	//m.Post("/games/new_request", handleGameRequest)
	//m.Post("/games/:id/register_server", handleRegisterServer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		fmt.Println(<-c)
		shutdown(machineId)
		os.Exit(1)
	}()
	defer shutdown(machineId)

	thisIp := fmt.Sprint("localhost:", strconv.Itoa(port))
	m.RunOnAddr(thisIp)
}

func handlePingRequest() (int, string) {
	return 200, "OK"
}

func shutdown(machineId string) {
	unregister := fmt.Sprintf("http://localhost:3000/machines/%s/unregister", machineId)
	_, err := http.PostForm(unregister, url.Values{"message": {"OK"}})
	if err != nil {
		fmt.Println(err)
	}

}

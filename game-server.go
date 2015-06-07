package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
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

	response, err := http.PostForm("http://localhost:3000/machines/register_new", url.Values{"port": {strconv.Itoa(port)}})
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

	m := martini.Classic()
	//m.Post("/games/new_request", handleGameRequest)
	//m.Post("/games/:id/register_server", handleRegisterServer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		shutdown(machineId)
		os.Exit(1)
	}()
	//defer shutdown(machineId)

	thisIp := fmt.Sprint("localhost:", strconv.Itoa(port))
	m.RunOnAddr(thisIp)

}

func shutdown(machineId string) {
	unregister := fmt.Sprintf("http://localhost:3000/machines/%s/unregister", machineId)
	_, err := http.PostForm(unregister, url.Values{"message": {"OK"}})
	if err != nil {
		fmt.Println(err)
	}

}

package main

import "fmt"
import "redis"
import "bytes"
import "os"

func main() {
	var input string
	fmt.Print("What game server chat u want in on: ")
	fmt.Scanln(&input)
	var channelKey string = "games:"+input+":chat"
	msg := connectClient(channelKey)
	go spamChat1(channelKey)
	chatbuf := bytes.NewBuffer(<-msg)
	chatbuf.WriteTo(os.Stdout)
	fmt.Scanln(&input)
		
}

func connectClient(chanKey string) (redis.PubSubChannel) {
	
	spec := redis.DefaultSpec().Password("go-redis")
	client, e := redis.NewPubSubClientWithSpec(spec)
	if e != nil {
		fmt.Print("Error creating client for: ", e)
	}
	defer client.Quit()
	fmt.Println("before subbing")
	client.Subscribe("chat")
	fmt.Println("after subbing")
	return client.Messages("chat")
	//fmt.Println(string(<-client.Messages("chat")))
	//fmt.Println(string(<-client.Messages("chat")))
	//for {	
	//	chatbuf := bytes.NewBuffer(<-client.Messages("chat"))
	//	chatbuf.WriteTo(os.Stdout)
	//	}
}

func spamChat1(chanKey string) {
	fmt.Println("spamchat being called") 
	spec := redis.DefaultSpec().Password("go-redis")
	client, e := redis.NewAsynchClientWithSpec(spec)
	if e != nil {
		fmt.Print("Error creating spam client for: ", e)
	}
	defer client.Quit()
	
	for {
		//var fr redis.FutureInt64
		//var fr2 redis.FutureBool
		var buf bytes.Buffer
		buf.Write([]byte("hello"))
		bt := buf.Bytes() 
		_, e = client.Publish("chat", bt)
		if e != nil {
			fmt.Println("error in publishing: ", e)
		}
		_, e = client.Rpush("chatlog",bt)
		if e != nil {
			fmt.Println("error in storing list: ", e)
		}
		
		//fr.Get()
		//fr2.Get()
	}
}




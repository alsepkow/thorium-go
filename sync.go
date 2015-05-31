package main

import "fmt"
import "time"


func worker(done chan bool, num int) {
	fmt.Print("starting work at ", num, "\n")
	
	time.Sleep(time.Duration(5)*time.Second);

	fmt.Print("finished work\n")

	done <- true
}


func main() {
	done := make(chan bool, 1)

	for i:=0; i<=5; i++ {
		go worker(done,i)
	}
	<- done;
}
	

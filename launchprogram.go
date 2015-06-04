package main

import "github.com/go-martini/martini"
import "os/exec"
import "fmt"
import "strconv"
import "redis"
import "bytes"
import "flag"
import "math/rand"
import "time"
import "database/sql"
import _ "github.com/lib/pq"



func Check(prog string) int {
	acceptedList := []string {"boltactiongame", "test"}
	for i:=0; i<len(acceptedList); i++ {
		if acceptedList[i] == prog {
			return 1
		}
	}
	return 0

}

//returns cmd struct of program
func Execute(prog string) *exec.Cmd {
	cmd := exec.Command("./"+prog)
	e := cmd.Start()
	if e!=nil {
		fmt.Println("Error runninng program ", e)
	}
	return cmd

}

func RedisPush(cmdI *exec.Cmd) int {
	spec := redis.DefaultSpec().Password("go-redis")
	client, e := redis.NewSynchClientWithSpec(spec)
	if e!= nil {
		fmt.Println("error creating client for: ", e)
	}
	defer client.Quit()
	pidString := strconv.Itoa(cmdI.Process.Pid)
	var buf bytes.Buffer
	buf.Write([]byte(pidString))
	e = client.Hset("server:pids", "pid", buf.Bytes())
	if e != nil {
		fmt.Println("error writing to list")
		return 0
	}
	return 1
}

func PostGresQueryIDS() []int {
	db, err := sql.Open("postgres", "user=thoriumnet password=thoriumtest dbname=thoriumnet host=localhost")
	if err != nil {
		fmt.Println("err: ", err)
	}
	var game_id int
	game_ids := make([]int, 100)
	rows, err := db.Query("SELECT * FROM games;")
	if err != nil {
		fmt.Println("err2: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&game_id)
		if err != nil {
			fmt.Println("err3: ", err)
		}
		for index,_ := range game_ids {
			if game_ids[index]==0 {
				game_ids[index]=game_id;
				break
			}
		}
	}
	return game_ids
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var portArg = flag.Int("p",rand.Intn(65000-10000)+10000, "specifies port, default is random int between 10000-65000")
	var mapArg = flag.String("m", "default map  value", "description of map")
	flag.Parse()
	fmt.Println(strconv.Itoa(*portArg))
	fmt.Println(*mapArg)
	processL := make([]*exec.Cmd, 100)
	currentGames := PostGresQueryIDS()
	for _,value := range currentGames {
		if value != 0 {
			fmt.Println(value)
		}
	}
	m := martini.Classic()
	m.Post("/launch/:name",  func(params martini.Params) string {
		e := Check(params["name"])
		if e==1 {
			cmdInfo := Execute(params["name"])
			for i:=0; i<len(processL); i++ {
				if processL[i]==nil {
					processL[i]=cmdInfo
					//suc := RedisPush(cmdInfo)
					RedisPush(cmdInfo)
					break
				}
			}
			//fmt.Println(processL)
			return "launching " + params["name"] + "with pid " + strconv.Itoa(cmdInfo.Process.Pid) 	
		} else {
			return "not accepted"
		}	
	})
	
	m.Run()
//	err := cmd.Wait()
//	fmt.Println(err)
//	fmt.Println(cmd.Path)
//	fmt.Println(cmd.Process.Pid)
}

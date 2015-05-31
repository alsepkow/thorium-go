package main

import "github.com/go-martini/martini"
import "os/exec"
import "fmt"
import "strconv"

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

func main() {
	//processL := [100]*exec.Cmd
	m := martini.Classic()
	m.Get("/launch/:name",  func(params martini.Params) string {
		e := Check(params["name"])
		if e==1 {
			cmdInfo := Execute(params["name"])
			//processList[0]=cmdInfo
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

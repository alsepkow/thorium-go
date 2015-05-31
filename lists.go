package main

import "fmt"
import "os"
import "os/exec"
import "time"

type ManagedProcess struct {
	ApplicationName string
	Process os.Process
	PID int
	// add more here if needed
}


func main(){
	fmt.Println("hello world")

	var n int
	n = 256

	var processes [n]ManagedProcess
	
	for i := 0; i < n; i++ {
		cmd := execute_cmd("./test")

		// todo: factor into managed process
		var process ManagedProcess
		process.Name = "test_outputonly"
		process.Process = cmd.Process
		process.PID = cmd.Process.Pid
	}
	
	fmt.Println("Thorium.NET is running ... ")
	
	time.Sleep(10);

	fmt.Println("Thorium.NET is shutting down ... ")
	
	for i = 0; i < n; i++ {
		processes[n].Process.Kill()
	}
}

func execute_cmd(command string) *exec.Cmd {
	cmd := exec.Command(command)
	err := cmd.Start()
	if (err != nil) {
		fmt.Println("Error: ", e)
		fmt.Println("User Commmand: ", command)
	}
	return cmd
}

// todo
//func startProcess(applicationName string) *ManagedProcess {

func addProcess(applicationName string) {
	cmd := e

}

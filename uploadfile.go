package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()
	m.Post("/upload/:fname", handleUploadFile)
	m.Run()
}

func handleUploadFile(req *http.Request, params martini.Params) (int, string) {
	//get current dir
	dir, e := filepath.Abs(filepath.Dir(""))
	if e != nil {
		fmt.Println("error with finding abs dir", e)
		return 500, "Internal server error"
	}

	fmt.Println(dir)
	fileName := params["fname"]

	//create file and name it based on url they passed
	outputFile, e2 := os.Create(dir + "/" + fileName)
	if e2 != nil {
		fmt.Println("can not create file", e)
		return 500, "Internal server error"
	}
	defer outputFile.Close()

	//multipart/form-data
	fmt.Println(req.Header["Content-Type"])

	//grab the stuff
	file, header_ptr, err := req.FormFile("filedata")
	if err != nil {
		fmt.Println("Error forming file: ", err)
		return 500, "Internal server error"
	}
	defer file.Close()

	fmt.Println(header_ptr.Filename)
	//fmt.Println(file)

	//copy stuff sent to us to our created file
	bytesWritten, err2 := io.Copy(outputFile, file)
	if err2 != nil {
		fmt.Println("unable to write to file", err2)
		return 500, "Internal server error"
	}

	//print how many bytes written (can compare to test, or can run program)
	fmt.Println(bytesWritten)
	return 200, fileName + " has been created on the server.\n"

}

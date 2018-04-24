package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger // Important information
	Error *log.Logger // Critical problem
)

const (
	errFile = "./sync_err.txt"
)

func init() {
	errFile, err := os.OpenFile(errFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(errFile, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func test() {
	Info.Println("This is info info...")
	Error.Println("This is err info...")
}

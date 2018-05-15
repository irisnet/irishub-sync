package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger // Important information
	Warning *log.Logger // Warning information
	Error *log.Logger // Critical problem
)

const (
	errFile = ".sync_server_err.log"
	warningFile = ".sync_server_warning.log"
)

func init() {
	errFile, err := os.OpenFile(os.ExpandEnv("$HOME/" + errFile),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	
	warningFile, err := os.OpenFile(os.ExpandEnv("$HOME/" + warningFile),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open warning log file:", err)
	}

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	
	Warning = log.New(io.MultiWriter(warningFile, os.Stderr),
		"Warning: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(errFile, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func test() {
	Info.Println("This is info info...")
	Warning.Println("This is warning info...")
	Error.Println("This is err info...")
}

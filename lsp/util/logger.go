package util

import (
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	f, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// do not close the file

	if err != nil {
		panic(err)
	}

	Logger = log.New(f, "lsp", log.LstdFlags|log.Lshortfile)
	Logger.Println("logger created")
}

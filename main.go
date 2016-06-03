package main

import (
	"./swearbot"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	LogFileName   = "log.txt"
	TokenFileName = "token.txt"
)

func main() {
	var logFile *os.File = createLogFile()
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	token := readSlackToken()
	swearbot.Run(token)
}

func createLogFile() *os.File {
	logFile, err := os.OpenFile(LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	return logFile
}

func readSlackToken() string {
	if _, err := os.Stat(TokenFileName); os.IsNotExist(err) {
		ioutil.WriteFile(TokenFileName, []byte("SLACK-TOKEN-HERE"), 0666)
		log.Fatalf("Enter slack token in %s", TokenFileName)
	}

	bytes, err := ioutil.ReadFile(TokenFileName)
	if err != nil {
		log.Fatalf("Cannot read slack token file '%s': %s", TokenFileName, err)
	}

	return strings.Trim(string(bytes), "\r\n ")
}

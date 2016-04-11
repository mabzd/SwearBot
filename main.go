package main

import (
	"./swearbot"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var logFile *os.File = createLogFile("log.txt")
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	token := readSlackToken("token.txt")
	config := readConfig("config.json")
	swearbot.Run(token, config)
}

func createLogFile(fileName string) *os.File {
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	return logFile
}

func readConfig(fileName string) swearbot.BotConfig {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Cannot read config from file '%s': %s", fileName, err)
	}

	var config swearbot.BotConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("Error when parsing config JSON: %s", err)
	}

	return config
}

func readSlackToken(fileName string) string {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		ioutil.WriteFile(fileName, []byte("SLACK-TOKEN-HERE"), 0666)
		log.Fatalf("Enter slack token in %s", fileName)
	}

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Cannot read slack token file '%s': %s", fileName, err)
	}

	return string(bytes)
}

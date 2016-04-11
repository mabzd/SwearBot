package main

import (
	"./swearbot"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Token     string
	BotConfig swearbot.BotConfig
}

func main() {
	var logFile *os.File = createLogFile("log.txt")
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	config := readConfig("config.json")
	bot := swearbot.NewSwearBot("swears.txt", "stats.json", config.BotConfig)
	bot.Run(config.Token)
}

func createLogFile(fileName string) *os.File {
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	return logFile
}

func readConfig(fileName string) Config {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Cannot read config from file '%s': %s", fileName, err)
	}
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("Error when parsing config JSON: %s", err)
	}
	return config
}

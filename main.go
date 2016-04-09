package main

import (
	"log"
	"os"
	"io"
	"io/ioutil"
	"encoding/json"
	"github.com/nlopes/slack"
	"./swearbot"
)

type Config struct {
	Token string
	BotConfig swearbot.BotConfig
}

func main() {
	var logFile *os.File = createLogFile("log.txt")
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	config := readConfig("config.json")

	api := slack.New(config.Token)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	log.Println("Start")
	processEvents(rtm, config)
	log.Println("Finish")
}

func createLogFile(fileName string) *os.File {
	logFile, err := os.OpenFile(fileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
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

func logInfo(info *slack.Info) {
	log.Println("Connected to: " + info.URL)
	for _, c := range info.Channels {
		if c.IsMember {
			log.Printf("Member of channel: #%s\n", c.Name)
		}
	}
}

func processEvents(rtm *slack.RTM, config Config) {
	swearBot := swearbot.NewSwearBot("swears.txt", config.BotConfig)
	swearBot.LoadSwears()
	log.Println("Swears loaded")

	for {
		select {
		case msg := <- rtm.IncomingEvents:

			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				logInfo(ev.Info)

			case *slack.MessageEvent:
				response := swearBot.ParseMessage(ev.Text)
				if response != "" {
					rtm.SendMessage(rtm.NewOutgoingMessage(response, ev.Channel))
				}

			case *slack.PresenceChangeEvent:
				// Ignore presence change

			case *slack.LatencyReport:
				// Ignore latency report

			case *slack.RTMError:
				log.Printf("RTM Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Println("Invalid credentials")
				return

			default:
				// Ignore other events
			}
		}
	}
}
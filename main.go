package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"io/ioutil"
	"strings"
	"github.com/nlopes/slack"
	"./swearbot"
)

func main() {
	var logFile *os.File = createLogFile("log.txt")
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	token := readBotToken("bot-token.txt")
	api := slack.New(token)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	log.Println("Start")
	processEvents(rtm)
	log.Println("Finish")
}

func createLogFile(fileName string) *os.File {
	logFile, err := os.OpenFile(fileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	return logFile
}

func readBotToken(fileName string) string {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Cannot read bot token from file '%s'", fileName)
	}
	return string(bytes)
}

func logInfo(info *slack.Info) {
	log.Println("Connected to: " + info.URL)
	for _, c := range info.Channels {
		if c.IsMember {
			log.Printf("Member of channel: #%s\n", c.Name)
		}
	}
}

func processEvents(rtm *slack.RTM) {
	swearBot := swearbot.NewSwearBot("swears.txt")
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
				swears := swearBot.FindSwears(ev.Text)
				if len(swears) > 0 {
					response := fmt.Sprintf("Following swears found: *%s*", strings.Join(swears, "*, *"))
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
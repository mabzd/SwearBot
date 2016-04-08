package main

import (
	"log"
	"os"
	"io"
	"io/ioutil"
	"bufio"
	"strings"
	"github.com/nlopes/slack"
	"./dictmatch"
)

func main() {
	var logFile *os.File = createLogFile("log.txt")
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	dict := dictmatch.NewDict()
	loadDict(dict, "swears.txt")
	log.Println("Swears loaded")

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

func loadDict(dict *dictmatch.Dict, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening swear dictionary file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := normalizeWord(scanner.Text())
		dict.AddEntry(word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from swar dictionary file: %v", err)
	}
}

func normalizeWord(word string) string {
	word = strings.Trim(word, " \n\r")
	word = strings.ToLower(word)
	return word
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
	for {
		select {
		case msg := <- rtm.IncomingEvents:

			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				logInfo(ev.Info)

			case *slack.MessageEvent:
				//rtm.SendMessage(rtm.NewOutgoingMessage("Test...", ev.Channel))

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
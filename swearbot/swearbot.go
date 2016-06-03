package swearbot

import (
	"../mods"
	"../mods/modswears"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"regexp"
)

var botMentionRegex *regexp.Regexp = nil

func registerMods() {

}

func Run(token string) {
	var connected bool = false

	slackClient := slack.New(token)
	slackClient.SetDebug(false)
	rtm := slackClient.NewRTM()

	modContainer := mods.NewModContainer()

	if !modContainer.LoadConfig() {
		log.Println("Loading mod config failed.")
		return
	}

	modContainer.AddMod(modswears.NewModSwears())

	if !modContainer.InitMods(slackClient) {
		log.Println("Initializing mods failed.")
		return
	}

	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:

			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				logInfo(ev.Info)
				compileMentionRegex(ev.Info.User.ID)
				connected = true

			case *slack.MessageEvent:
				if connected {
					response := ""
					message := ev.Text
					userId := ev.User
					channel := ev.Channel

					if isMention(message) {
						message = removeMentions(message)
						response = modContainer.ProcessMention(message, userId, channel)
					} else {
						response = modContainer.ProcessMessage(message, userId, channel)
					}

					respond(rtm, response, channel)
				}

			case *slack.RTMError:
				log.Printf("RTM Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Println("Invalid credentials")
				return

			default:
			}
		}
	}
}

func compileMentionRegex(botId string) {
	expr := fmt.Sprintf("<@%s>:?", regexp.QuoteMeta(botId))
	botMentionRegex = regexp.MustCompile(expr)
}

func isMention(message string) bool {
	return botMentionRegex.MatchString(message)
}

func removeMentions(message string) string {
	return botMentionRegex.ReplaceAllLiteralString(message, "")
}

func respond(rtm *slack.RTM, response string, channel string) {
	if response != "" {
		rtm.SendMessage(rtm.NewOutgoingMessage(response, channel))
	}
}

func logInfo(info *slack.Info) {
	log.Println("Connected to: " + info.URL)
	log.Printf("Bot name: @%s", info.User.Name)
	for _, c := range info.Channels {
		if c.IsMember {
			log.Printf("Member of channel: #%s\n", c.Name)
		}
	}
}

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
var connected bool = false

func Run(token string) {
	slackClient := slack.New(token)
	slackClient.SetDebug(false)
	rtm := slackClient.NewRTM()
	modContainer := createModContainer(slackClient)
	if modContainer == nil {
		log.Println("Creating mods failed.")
		return
	}
	go rtm.ManageConnection()
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch event := msg.Data.(type) {
			case *slack.ConnectedEvent:
				onConnect(event.Info)
			case *slack.MessageEvent:
				onMessage(rtm, event, modContainer)
			case *slack.RTMError:
				onError(event)
			case *slack.InvalidAuthEvent:
				log.Println("Invalid credentials")
				return
			default:
			}
		}
	}
}

func createModContainer(slackClient *slack.Client) *mods.ModContainer {
	modContainer := mods.NewModContainer()
	if !modContainer.LoadConfig() {
		return nil
	}
	modContainer.AddMod(modswears.NewModSwears())
	if !modContainer.InitMods(slackClient) {
		return nil
	}
	return modContainer
}

func onConnect(info *slack.Info) {
	logInfo(info)
	compileMentionRegex(info.User.ID)
	connected = true
}

func onMessage(
	rtm *slack.RTM,
	event *slack.MessageEvent,
	modContainer *mods.ModContainer) {

	if connected {
		response := ""
		message := event.Text
		userId := event.User
		channelId := event.Channel
		if isMention(message) {
			message = removeMentions(message)
			response = modContainer.ProcessMention(message, userId, channelId)
		} else {
			response = modContainer.ProcessMessage(message, userId, channelId)
		}
		respond(rtm, response, channelId)
	}
}

func onError(err *slack.RTMError) {
	log.Printf("RTM Error: %s\n", err.Error())
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

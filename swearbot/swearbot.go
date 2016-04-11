package swearbot

import (
	"../swears"
	"github.com/nlopes/slack"
	"log"
)

type BotConfig struct {
	SwearsConfig swears.SwearsConfig
}

func Run(token string, config BotConfig) {

	api := slack.New(token)
	api.SetDebug(false)
	rtm := api.NewRTM()

	swears := swears.NewSwears(api, config.SwearsConfig)
	swears.LoadSwears()

	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:

			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				logInfo(ev.Info)

			case *slack.MessageEvent:
				response := swears.ProcessMessage(ev.Text, ev.User)
				if response != "" {
					rtm.SendMessage(rtm.NewOutgoingMessage(response, ev.Channel))
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

func logInfo(info *slack.Info) {
	log.Println("Connected to: " + info.URL)
	log.Printf("Bot name: @%s", info.User.Name)
	for _, c := range info.Channels {
		if c.IsMember {
			log.Printf("Member of channel: #%s\n", c.Name)
		}
	}
}

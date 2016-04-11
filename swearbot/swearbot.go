package swearbot

import (
	"../stats"
	"../swears"
	"bytes"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"regexp"
	"strings"
	"time"
)

type BotConfig struct {
	AddRuleRegex          string
	MonthlyRankRegex      string
	OnAddRuleResponse     string
	OnSwearsFoundResponse string
	OnEmptyRankResponse   string
	OnUserFetchErr        string
	SwearsConfig          swears.SwearsConfig
	StatsConfig           stats.StatsConfig
}

type SwearBot struct {
	name             string
	api              *slack.Client
	swears           *swears.Swears
	stats            *stats.Stats
	addRuleRegex     *regexp.Regexp
	monthlyRankRegex *regexp.Regexp
	config           BotConfig
}

func NewSwearBot(
	dictFileName string,
	statsFileName string,
	config BotConfig) *SwearBot {

	return &SwearBot{
		name:             "",
		api:              nil,
		swears:           swears.NewSwears(dictFileName, config.SwearsConfig),
		stats:            stats.NewStats(statsFileName, config.StatsConfig),
		addRuleRegex:     regexp.MustCompile(config.AddRuleRegex),
		monthlyRankRegex: regexp.MustCompile(config.MonthlyRankRegex),
		config:           config,
	}
}

func (sb *SwearBot) Run(token string) {

	sb.swears.LoadSwears()
	sb.api = slack.New(token)
	sb.api.SetDebug(false)
	rtm := sb.api.NewRTM()

	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:

			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				logInfo(ev.Info)
				sb.name = ev.Info.User.Name

			case *slack.MessageEvent:
				response := sb.ParseMessage(ev.Text, ev.User)
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

func (sb *SwearBot) ParseMessage(message string, user string) string {
	if sb.monthlyRankRegex.MatchString(message) {
		return sb.printMonthlyRank()
	}

	rules := sb.addRuleRegex.FindAllStringSubmatch(message, 1)
	if rules != nil {
		return sb.addRule(rules[0][1])
	}

	return sb.parseSwears(message, user)
}

func (sb *SwearBot) printMonthlyRank() string {
	now := time.Now()
	userStats, rankErr := sb.stats.GetMonthlyRank(int(now.Month()), now.Year())
	if rankErr != nil {
		return rankErr.Error()
	}

	if len(userStats) == 0 {
		return sb.config.OnEmptyRankResponse
	}

	users, usersErr := sb.api.GetUsers()
	if usersErr != nil {
		log.Printf("Print monthly rank: cannot fetch users from slack: %s\n", usersErr)
		return sb.config.OnUserFetchErr
	}

	var response bytes.Buffer
	for i, userStat := range userStats {
		user := getUserById(users, userStat.UserId)
		rankLine := fmt.Sprintf("%d. *%s*: %d swears\n", i+1, user.Name, userStat.SwearCount)
		response.WriteString(rankLine)
	}

	return response.String()
}

func (sb *SwearBot) addRule(rule string) string {
	err := sb.swears.AddRule(rule)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf(sb.config.OnAddRuleResponse, rule)
}

func (sb *SwearBot) parseSwears(message string, user string) string {
	swears := sb.swears.FindSwears(message)

	if len(swears) > 0 {
		now := time.Now()
		err := sb.stats.AddSwearCount(int(now.Month()), now.Year(), user, len(swears))
		if err != nil {
			return err.Error()
		}

		swearsLine := fmt.Sprintf("*%s*", strings.Join(swears, "*, *"))
		response := fmt.Sprintf(sb.config.OnSwearsFoundResponse, swearsLine)
		return response
	}

	return ""
}

func getUserById(users []slack.User, id string) *slack.User {
	for _, user := range users {
		if user.ID == id {
			return &user
		}
	}
	return nil
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

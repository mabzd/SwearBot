package swears

import (
	"../dictmatch"
	"bytes"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"regexp"
	"strings"
	"time"
)

const Success = 0

type Swears struct {
	api              *slack.Client
	dict             *dictmatch.Dict
	addRuleRegex     *regexp.Regexp
	monthlyRankRegex *regexp.Regexp
	config           SwearsConfig
}

type SwearsConfig struct {
	DictFileName  string
	StatsFileName string

	AddRuleRegex     string
	MonthlyRankRegex string

	OnUnknownCommandResponse string
	OnAddRuleResponse        string
	OnSwearsFoundResponse    string
	OnEmptyRankResponse      string

	OnUserFetchErr       string
	OnAddRuleFileReadErr string
	OnAddRuleConflictErr string
	OnAddRuleSaveErr     string
	OnInvalidWildcardErr string

	OnStatsFileCreateErr string
	OnStatsFileReadErr   string
	OnStatsUnmarshalErr  string
	OnStatsMarshalErr    string
	OnStatsSaveErr       string
}

func NewSwears(api *slack.Client, config SwearsConfig) *Swears {
	return &Swears{
		api:              api,
		dict:             dictmatch.NewDict(),
		addRuleRegex:     regexp.MustCompile(config.AddRuleRegex),
		monthlyRankRegex: regexp.MustCompile(config.MonthlyRankRegex),
		config:           config,
	}
}

func (sw *Swears) ProcessMention(message string, userId string) string {
	if sw.monthlyRankRegex.MatchString(message) {
		return sw.printMonthlyRank()
	}

	rules := sw.addRuleRegex.FindAllStringSubmatch(message, 1)
	if rules != nil {
		return sw.addRule(rules[0][1])
	}

	return sw.config.OnUnknownCommandResponse
}

func (sw *Swears) ProcessMessage(message string, userId string) string {
	swears := sw.FindSwears(message)

	if len(swears) > 0 {
		now := time.Now()
		err := sw.AddSwearCount(int(now.Month()), now.Year(), userId, len(swears))
		if err != Success {
			return getResponseOnErr(err, sw.config)
		}

		swearsLine := fmt.Sprintf("*%s*", strings.Join(swears, "*, *"))
		response := fmt.Sprintf(sw.config.OnSwearsFoundResponse, swearsLine)
		return response
	}

	return ""
}

func (sw *Swears) printMonthlyRank() string {
	now := time.Now()
	userStats, rankErr := sw.GetMonthlyRank(int(now.Month()), now.Year())
	if rankErr != Success {
		return getResponseOnErr(rankErr, sw.config)
	}

	if len(userStats) == 0 {
		return sw.config.OnEmptyRankResponse
	}

	users, usersErr := sw.api.GetUsers()
	if usersErr != nil {
		log.Printf("Monthly rank: Cannot fetch users from slack: %s\n", usersErr)
		return sw.config.OnUserFetchErr
	}

	var response bytes.Buffer
	for i, userStat := range userStats {
		user := getUserById(users, userStat.UserId)
		rankLine := fmt.Sprintf("%d. *%s*: %d swears\n", i+1, user.Name, userStat.SwearCount)
		response.WriteString(rankLine)
	}

	return response.String()
}

func (sw *Swears) addRule(rule string) string {
	err := sw.AddRule(rule)
	if err != Success {
		return getResponseOnErr(err, sw.config)
	}

	return fmt.Sprintf(sw.config.OnAddRuleResponse, rule)
}

func getUserById(users []slack.User, id string) *slack.User {
	for _, user := range users {
		if user.ID == id {
			return &user
		}
	}
	return nil
}

func getResponseOnErr(err int, config SwearsConfig) string {
	switch err {
	case AddRuleFileReadErr:
		return config.OnAddRuleFileReadErr
	case AddRuleConflictErr:
		return config.OnAddRuleConflictErr
	case AddRuleSaveErr:
		return config.OnAddRuleSaveErr
	case InvalidWildcardErr:
		return config.OnInvalidWildcardErr
	case StatsFileCreateErr:
		return config.OnStatsFileCreateErr
	case StatsFileReadErr:
		return config.OnStatsFileReadErr
	case StatsUnmarshalErr:
		return config.OnStatsUnmarshalErr
	case StatsMarshalErr:
		return config.OnStatsMarshalErr
	case StatsSaveErr:
		return config.OnStatsSaveErr
	default:
		log.Printf("Swears: No response for error code %d!", err)
		return fmt.Sprintf("Error #%v", err)
	}
}

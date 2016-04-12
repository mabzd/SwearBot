package swears

import (
	"../dictmatch"
	"../utils"
	"bytes"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"regexp"
	"strconv"
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
	DictFileName     string
	StatsFileName    string
	SettingsFileName string

	AddRuleRegex     string
	MonthlyRankRegex string

	SwearFormat              string
	OnSwearsFoundResponse    string
	OnUnknownCommandResponse string
	OnAddRuleResponse        string
	OnEmptyRankResponse      string
	MonhlyRankHeaderFormat   string
	RankLineFormat           string
	MonthNames               []string

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

		return formatSwearsResponse(sw.config.OnSwearsFoundResponse, sw.config.SwearFormat, swears)
	}

	return ""
}

func (sw *Swears) printMonthlyRank() string {
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	userStats, rankErr := sw.GetMonthlyRank(month, year)
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

	return formatMonthlyRank(sw.config, month, year, users, userStats)
}

func (sw *Swears) addRule(rule string) string {
	err := sw.AddRule(rule)
	if err != Success {
		return getResponseOnErr(err, sw.config)
	}

	return formatAddRuleResponse(sw.config.OnAddRuleResponse, rule)
}

func getUserById(users []slack.User, id string) (slack.User, bool) {
	for _, user := range users {
		if user.ID == id {
			return user, true
		}
	}
	return slack.User{}, false
}

func formatAddRuleResponse(format string, rule string) string {
	params := map[string]string{"rule": rule}
	return utils.ParamFormat(format, params)
}

func formatSwearsResponse(lineFormat string, swearFormat string, swears []string) string {
	var buffer bytes.Buffer
	for i, swear := range swears {
		buffer.WriteString(formatSwear(swearFormat, swear, i+1))
		buffer.WriteString(", ")
	}

	result := strings.Trim(buffer.String(), ", ")
	params := map[string]string{"swears": result, "count": strconv.Itoa(len(swears))}
	return utils.ParamFormat(lineFormat, params)
}

func formatSwear(format string, swear string, index int) string {
	params := map[string]string{"swear": swear, "index": strconv.Itoa(index)}
	return utils.ParamFormat(format, params)
}

func formatMonthlyRank(
	config SwearsConfig,
	month int,
	year int,
	users []slack.User,
	userStats []*UserStats) string {

	header := formatMonthlyRankHeader(config.MonhlyRankHeaderFormat, config.MonthNames, month, year)
	rankLines := formatRankLines(config.RankLineFormat, users, userStats)
	return fmt.Sprintf("%s\n%s", header, rankLines)
}

func formatMonthlyRankHeader(headerFormat string, monthNames []string, month int, year int) string {
	params := map[string]string{
		"month":    monthNames[month-1],
		"monthnum": strconv.Itoa(month),
		"year":     strconv.Itoa(year),
	}
	return utils.ParamFormat(headerFormat, params)
}

func formatRankLines(lineFormat string, users []slack.User, userStats []*UserStats) string {
	var buffer bytes.Buffer
	for i, userStat := range userStats {
		user, ok := getUserById(users, userStat.UserId)
		if !ok {
			user.Name = "unknown"
		}

		buffer.WriteString(formatRankLine(lineFormat, user, userStat.SwearCount, i+1))
		buffer.WriteString("\n")
	}

	return buffer.String()
}

func formatRankLine(lineFormat string, user slack.User, count int, index int) string {
	params := map[string]string{
		"index": strconv.Itoa(index),
		"user":  user.Name,
		"count": strconv.Itoa(count),
	}
	return utils.ParamFormat(lineFormat, params)
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

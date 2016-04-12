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
	api                 *slack.Client
	dict                *dictmatch.Dict
	addRuleRegex        *regexp.Regexp
	monthlyRankRegex    *regexp.Regexp
	swearNotifyOnRegex  *regexp.Regexp
	swearNotifyOffRegex *regexp.Regexp
	settings            *AllSettings
	config              SwearsConfig
}

type SwearsConfig struct {
	DictFileName     string
	StatsFileName    string
	SettingsFileName string

	AddRuleRegex        string
	MonthlyRankRegex    string
	SwearNotifyOnRegex  string
	SwearNotifyOffRegex string

	SwearFormat              string
	OnSwearsFoundResponse    string
	OnUnknownCommandResponse string
	OnAddRuleResponse        string
	OnEmptyRankResponse      string
	OnSwearNotifyOnResponse  string
	OnSwearNotifyOffResponse string
	MonhlyRankHeaderFormat   string
	RankLineFormat           string
	MonthNames               []string

	OnUserFetchErr       string
	OnDictFileReadErr    string
	OnAddRuleConflictErr string
	OnAddRuleSaveErr     string
	OnInvalidWildcardErr string

	OnStatsFileCreateErr string
	OnStatsFileReadErr   string
	OnStatsUnmarshalErr  string
	OnStatsMarshalErr    string
	OnStatsSaveErr       string

	OnSettingsFileCreateErr string
	OnSettingsFileReadErr   string
	OnSettingsUnmarshalErr  string
	OnSettingsMarshalErr    string
	OnSettingsSaveErr       string
}

func NewSwears(api *slack.Client, config SwearsConfig) *Swears {
	return &Swears{
		api:      api,
		dict:     dictmatch.NewDict(),
		settings: &AllSettings{UserSettings: map[string][]*UserSettings{}},
		config:   config,
	}
}

func (sw *Swears) Init() bool {
	var err error
	var errnum int

	sw.addRuleRegex, err = regexp.Compile(sw.config.AddRuleRegex)
	if err != nil {
		log.Printf("Swears: cannot compile AddRuleRegex: %v", err)
		return false
	}

	sw.monthlyRankRegex, err = regexp.Compile(sw.config.MonthlyRankRegex)
	if err != nil {
		log.Printf("Swears: cannot compile MonthlyRankRegex: %v", err)
		return false
	}

	sw.swearNotifyOnRegex, err = regexp.Compile(sw.config.SwearNotifyOnRegex)
	if err != nil {
		log.Printf("Swears: cannot compile SwearNotifyOnRegex: %v", err)
		return false
	}

	sw.swearNotifyOffRegex, err = regexp.Compile(sw.config.SwearNotifyOffRegex)
	if err != nil {
		log.Printf("Swears: cannot compile SwearNotifyOffRegex: %v", err)
		return false
	}

	errnum = sw.LoadSwears()
	if errnum != Success {
		return false
	}

	errnum = sw.LoadSettings()
	if errnum != Success {
		return false
	}

	return true
}

func (sw *Swears) ProcessMention(message string, userId string, channel string) string {
	if sw.monthlyRankRegex.MatchString(message) {
		return sw.printMonthlyRank()
	}

	rules := sw.addRuleRegex.FindAllStringSubmatch(message, 1)
	if rules != nil {
		return sw.addRule(rules[0][1])
	}

	if sw.swearNotifyOnRegex.MatchString(message) {
		return sw.setSwearNotify(userId, channel, "on")
	}

	if sw.swearNotifyOffRegex.MatchString(message) {
		return sw.setSwearNotify(userId, channel, "off")
	}

	return sw.config.OnUnknownCommandResponse
}

func (sw *Swears) ProcessMessage(message string, userId string, channel string) string {
	swears := sw.FindSwears(message)

	if len(swears) > 0 {
		now := time.Now()
		err := sw.AddSwearCount(int(now.Month()), now.Year(), userId, len(swears))
		if err != Success {
			return getResponseOnErr(err, sw.config)
		}

		swearNotify, exist := sw.settings.GetSetting(userId, channel, "SwearNotify")
		if exist && swearNotify == "on" {
			return formatSwearsResponse(
				sw.config.OnSwearsFoundResponse,
				sw.config.SwearFormat,
				swears)
		}
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

func (sw *Swears) setSwearNotify(userId string, channel string, value string) string {
	sw.settings.SetSetting(userId, channel, "SwearNotify", value)
	err := sw.SaveSettings()
	if err != Success {
		return getResponseOnErr(err, sw.config)
	}

	if value == "on" {
		return sw.config.OnSwearNotifyOnResponse
	}

	return sw.config.OnSwearNotifyOffResponse
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
	case DictFileReadErr:
		return config.OnDictFileReadErr
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
	case SettingsFileCreateErr:
		return config.OnSettingsFileCreateErr
	case SettingsFileReadErr:
		return config.OnSettingsFileReadErr
	case SettingsUnmarshalErr:
		return config.OnSettingsUnmarshalErr
	case SettingsMarshalErr:
		return config.OnSettingsMarshalErr
	case SettingsSaveErr:
		return config.OnSettingsSaveErr
	default:
		log.Printf("Swears: No response for error code %d!", err)
		return fmt.Sprintf("Error #%v", err)
	}
}

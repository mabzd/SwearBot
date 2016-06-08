package modswears

import (
	"../../dictmatch"
	"../../mods"
	"../../settings"
	"../../utils"
	"bytes"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	Success        = 0
	ConfigFileName = "config.json"
	DictFileName   = "swears.txt"
	StatsFileName  = "stats.json"
)

const (
	SettingSwearNotify = "ModSwears.SwearNotify"
)

type ModSwears struct {
	state               *mods.ModState
	dict                *dictmatch.Dict
	addRuleRegex        *regexp.Regexp
	currMonthRankRegex  *regexp.Regexp
	prevMonthRankRegex  *regexp.Regexp
	totalRankRegex      *regexp.Regexp
	swearNotifyOnRegex  *regexp.Regexp
	swearNotifyOffRegex *regexp.Regexp
	config              *ModSwearsConfig
	dictFileName        string
	statsFileName       string
}

func NewModSwears() *ModSwears {
	return &ModSwears{
		dict:   dictmatch.NewDict(),
		config: NewModSwearsConfig(),
	}
}

func (mod *ModSwears) Name() string {
	return "modswears"
}

func (mod *ModSwears) Init(state *mods.ModState) bool {
	var err error
	var errnum int
	mod.state = state
	mod.dictFileName = mods.GetPath(mod, DictFileName)
	mod.statsFileName = mods.GetPath(mod, StatsFileName)
	configFileName := mods.GetPath(mod, ConfigFileName)
	err = utils.JsonFromFileCreate(configFileName, mod.config)
	if err != nil {
		log.Println("ModSwears: cannot load config.")
		return false
	}
	mod.addRuleRegex, err = regexp.Compile(mod.config.AddRuleRegex)
	if err != nil {
		log.Printf("ModSwears: cannot compile AddRuleRegex: %v\n", err)
		return false
	}
	mod.currMonthRankRegex, err = regexp.Compile(mod.config.CurrMonthRankRegex)
	if err != nil {
		log.Printf("ModSwears: cannot compile CurrMonthRankRegex: %v\n", err)
		return false
	}
	mod.prevMonthRankRegex, err = regexp.Compile(mod.config.PrevMonthRankRegex)
	if err != nil {
		log.Printf("ModSwears: cannot compile PrevMonthRankRegex: %v\n", err)
	}
	mod.totalRankRegex, err = regexp.Compile(mod.config.TotalRankRegex)
	if err != nil {
		log.Printf("ModSwears: cannot compile TotalRankRegex: %v\n", err)
	}
	mod.swearNotifyOnRegex, err = regexp.Compile(mod.config.SwearNotifyOnRegex)
	if err != nil {
		log.Printf("ModSwears: cannot compile SwearNotifyOnRegex: %v\n", err)
		return false
	}
	mod.swearNotifyOffRegex, err = regexp.Compile(mod.config.SwearNotifyOffRegex)
	if err != nil {
		log.Printf("ModSwears: cannot compile SwearNotifyOffRegex: %v\n", err)
		return false
	}
	errnum = mod.LoadSwears()
	if errnum != Success {
		log.Println("ModSwears: loading swears dictionary failed.")
		return false
	}
	return true
}

func (mod *ModSwears) ProcessMention(message string, userId string, channelId string) string {
	if mod.currMonthRankRegex.MatchString(message) {
		return mod.getCurrMonthRank()
	}
	if mod.prevMonthRankRegex.MatchString(message) {
		return mod.getPrevMonthRank()
	}
	if mod.totalRankRegex.MatchString(message) {
		return mod.getTotalRank()
	}
	rules := mod.addRuleRegex.FindAllStringSubmatch(message, 1)
	if rules != nil {
		return mod.addRule(rules[0][1])
	}
	if mod.swearNotifyOnRegex.MatchString(message) {
		return mod.setSwearNotify(userId, channelId, "on")
	}
	if mod.swearNotifyOffRegex.MatchString(message) {
		return mod.setSwearNotify(userId, channelId, "off")
	}
	return ""
}

func (mod *ModSwears) ProcessMessage(message string, userId string, channelId string) string {
	swears := mod.FindSwears(message)
	if len(swears) > 0 {
		now := utils.TimeClock.Now()
		err := mod.AddSwearCount(int(now.Month()), now.Year(), userId, len(swears))
		if err != Success {
			return getResponseOnErr(err, mod.config)
		}
		swearNotify, exist := mod.state.GetUserChanSetting(
			userId,
			channelId,
			SettingSwearNotify)
		if exist && swearNotify == "on" {
			return formatSwearsResponse(
				mod.config.OnSwearsFoundResponse,
				mod.config.SwearFormat,
				swears)
		}
	}
	return ""
}

func (mod *ModSwears) getCurrMonthRank() string {
	now := utils.TimeClock.Now()
	month := int(now.Month())
	year := now.Year()
	return mod.getRankByMonth(month, year)
}

func (mod *ModSwears) getPrevMonthRank() string {
	prevMonth := utils.LastDayOfPrevMonth(utils.TimeClock.Now())
	month := int(prevMonth.Month())
	year := prevMonth.Year()
	return mod.getRankByMonth(month, year)
}

func (mod *ModSwears) getTotalRank() string {
	userStats, rankErr := mod.GetTotalRank()
	response := mod.prepareRank(userStats, rankErr)
	if response != "" {
		return response
	}
	return formatTotalRank(mod.config, userStats)
}

func (mod *ModSwears) getRankByMonth(month int, year int) string {
	userStats, rankErr := mod.GetMonthlyRank(month, year)
	response := mod.prepareRank(userStats, rankErr)
	if response != "" {
		return response
	}
	return formatMonthlyRank(mod.config, month, year, userStats)
}

func (mod *ModSwears) prepareRank(userStats []*UserStats, rankErr int) string {
	if rankErr != Success {
		return getResponseOnErr(rankErr, mod.config)
	}
	if len(userStats) == 0 {
		return mod.config.OnEmptyRankResponse
	}
	return fillUserRealNames(userStats, mod.state.SlackClient, mod.config)
}

func (mod *ModSwears) addRule(rule string) string {
	err := mod.AddRule(rule)
	if err != Success {
		return getResponseOnErr(err, mod.config)
	}
	return formatAddRuleResponse(mod.config.OnAddRuleResponse, rule)
}

func (mod *ModSwears) setSwearNotify(
	userId string,
	channelId string,
	value string) string {

	mod.state.SetUserChanSetting(userId, channelId, SettingSwearNotify, value)
	err := mod.state.Save()
	if err != Success {
		return getResponseOnErr(err, mod.config)
	}
	if value == "on" {
		return mod.config.OnSwearNotifyOnResponse
	}
	return mod.config.OnSwearNotifyOffResponse
}

func fillUserRealNames(
	userStats []*UserStats,
	slack *slack.Client,
	config *ModSwearsConfig) string {

	users, usersErr := slack.GetUsers()
	if usersErr != nil {
		log.Printf("ModSwears: Cannot fetch users from slack: %s\n", usersErr)
		return config.OnUserFetchErr
	}
	for _, userStat := range userStats {
		user, ok := getUserById(users, userStat.UserId)
		if !ok {
			userStat.UserId = "unknown"
		} else {
			userStat.UserId = user.Name
		}
	}
	return ""
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

func formatSwearsResponse(
	lineFormat string,
	swearFormat string,
	swears []string) string {

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
	config *ModSwearsConfig,
	month int,
	year int,
	userStats []*UserStats) string {

	header := formatMonthlyRankHeader(
		config.MonthlyRankHeaderFormat,
		config.MonthNames,
		month,
		year)
	rankLines := formatRankLines(config.RankLineFormat, userStats)
	return fmt.Sprintf("%s\n%s", header, rankLines)
}

func formatTotalRank(
	config *ModSwearsConfig,
	userStats []*UserStats) string {

	header := config.TotalRankHeaderFormat
	rankLines := formatRankLines(config.RankLineFormat, userStats)
	return fmt.Sprintf("%s\n%s", header, rankLines)
}

func formatMonthlyRankHeader(
	headerFormat string,
	monthNames []string,
	month int,
	year int) string {

	params := map[string]string{
		"month":    monthNames[month-1],
		"monthnum": strconv.Itoa(month),
		"year":     strconv.Itoa(year),
	}
	return utils.ParamFormat(headerFormat, params)
}

func formatRankLines(lineFormat string, userStats []*UserStats) string {
	var buffer bytes.Buffer
	for i, userStat := range userStats {
		line := formatRankLine(lineFormat, userStat.UserId, userStat.SwearCount, i+1)
		buffer.WriteString(line)
		buffer.WriteString("\n")
	}

	return buffer.String()
}

func formatRankLine(lineFormat string, user string, count int, index int) string {
	params := map[string]string{
		"index": strconv.Itoa(index),
		"user":  user,
		"count": strconv.Itoa(count),
	}
	return utils.ParamFormat(lineFormat, params)
}

//TODO: move settings responses to some general config
func getResponseOnErr(err int, config *ModSwearsConfig) string {
	switch err {
	case DictFileReadErr:
		return config.OnDictFileReadErr
	case AddRuleConflictErr:
		return config.OnAddRuleConflictErr
	case AddRuleSaveErr:
		return config.OnAddRuleSaveErr
	case InvalidWildcardErr:
		return config.OnInvalidWildcardErr
	case StatsFileReadErr:
		return config.OnStatsFileReadErr
	case StatsSaveErr:
		return config.OnStatsSaveErr
	case settings.SettingsFileReadErr:
		return config.OnSettingsFileReadErr
	case settings.SettingsSaveErr:
		return config.OnSettingsSaveErr
	default:
		log.Printf("ModSwears: No response for error code %d!\n", err)
		return fmt.Sprintf("Error #%v", err)
	}
}

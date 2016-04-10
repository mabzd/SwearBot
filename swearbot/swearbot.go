package swearbot

import (
	"../stats"
	"../swears"
	"fmt"
	"regexp"
	"strings"
)

type BotConfig struct {
	AddRuleRegex          string
	OnAddRuleResponse     string
	OnSwearsFoundResponse string
	SwearsConfig          swears.SwearsConfig
	StatsConfig           stats.StatsConfig
}

type SwearBot struct {
	swears       *swears.Swears
	stats        *stats.Stats
	addRuleRegex *regexp.Regexp
	config       BotConfig
}

func NewSwearBot(dictFileName string, statsFileName string, config BotConfig) *SwearBot {
	return &SwearBot{
		swears:       swears.NewSwears(dictFileName, config.SwearsConfig),
		stats:        stats.NewStats(statsFileName, config.StatsConfig),
		addRuleRegex: regexp.MustCompile(config.AddRuleRegex),
		config:       config,
	}
}

func (sb *SwearBot) LoadSwears() {
	sb.swears.LoadSwears()
}

func (sb *SwearBot) ParseMessage(message string) string {
	rules := sb.addRuleRegex.FindAllStringSubmatch(message, 1)
	if rules != nil {
		rule := rules[0][1]
		err := sb.swears.AddRule(rule)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf(sb.config.OnAddRuleResponse, rule)
	}
	return sb.parseSwears(message)
}

func (sb *SwearBot) parseSwears(message string) string {
	swears := sb.swears.FindSwears(message)
	if len(swears) > 0 {
		swearsLine := fmt.Sprintf("*%s*", strings.Join(swears, "*, *"))
		response := fmt.Sprintf(sb.config.OnSwearsFoundResponse, swearsLine)
		return response
	}
	return ""
}

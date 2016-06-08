package modchoice

import (
	"../../mods"
	"../../utils"
	"log"
	"regexp"
	"strings"
)

const (
	ConfigFileName = "config.json"
)

var wordSplitRegex *regexp.Regexp = regexp.MustCompile("\\s+")

type ModChoice struct {
	state  mods.State
	config *ModChoiceConfig
}

func NewModChoice() *ModChoice {
	return &ModChoice{
		config: NewModChoiceConfig(),
	}
}

func (mod *ModChoice) Name() string {
	return "modchoice"
}

func (mod *ModChoice) Init(state mods.State) bool {
	var err error
	mod.state = state
	configFilePath := mods.GetPath(mod, ConfigFileName)
	err = utils.JsonFromFileCreate(configFilePath, mod.config)
	if err != nil {
		log.Printf("ModChoice: cannot load config")
		return false
	}
	return validateConfig(mod.config)
}

func (mod *ModChoice) ProcessMention(
	message string,
	userId string,
	channelId string) string {

	options := getOptions(message, mod.config.OrKeywords)
	if len(options) > 1 {
		return mod.choose(options)
	}
	return ""
}

func (mod *ModChoice) ProcessMessage(
	message string,
	userId string,
	channelId string) string {

	return ""
}

func (mod *ModChoice) choose(options []string) string {
	if utils.RandEvent(mod.config.NullChoiceProbability) {
		return utils.RandSelect(mod.config.NullChoiceResponses)
	}
	option := utils.RandSelect(options)
	response := utils.RandSelect(mod.config.ChoiceResponseFormat)
	return formatResponse(response, option)
}

func getOptions(message string, orKeys []string) []string {
	message = strings.TrimRight(message, " ?\r\n\t")
	words := getWords(message)
	options := []string{}
	j := 0
	for i, word := range words {
		if utils.ContainsCaseIns(word, orKeys) {
			option := strings.Join(words[j:i], " ")
			j = i + 1
			if len(option) > 0 {
				options = append(options, option)
			}
		}
	}
	if j < len(words) {
		options = append(options, strings.Join(words[j:], " "))
	}
	return options
}

func getWords(message string) []string {
	words := wordSplitRegex.Split(message, -1)
	result := []string{}
	for _, word := range words {
		word = strings.Trim(word, " \r\n\t")
		if len(word) > 0 {
			result = append(result, word)
		}
	}
	return result
}

func validateConfig(config *ModChoiceConfig) bool {
	for _, orKey := range config.OrKeywords {
		if len(getWords(orKey)) > 1 {
			log.Printf("ModChoice: OR keyword '%s' must be a single word.", orKey)
			return false
		}
	}
	return true
}

func formatResponse(format string, option string) string {
	params := map[string]string{"option": option}
	return utils.ParamFormat(format, params)
}

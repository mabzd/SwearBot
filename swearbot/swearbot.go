package swearbot

import (
	"errors"
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"regexp"
	"../dictmatch"
)

type BotConfig struct {
	AddRuleRegex string
	OnSwearsFoundResponse string
	OnAddRuleResponse string
	OnAddRuleFileReadErr string
	OnAddRuleConflictErr string
	OnAddRuleSaveErr string
}

type SwearBot struct {
	dict *dictmatch.Dict
	dictFileName string
	addRuleRegex *regexp.Regexp
	config BotConfig
}

func NewSwearBot(fileName string, botConfig BotConfig) *SwearBot {
	return &SwearBot {
		dict: dictmatch.NewDict(),
		dictFileName: fileName,
		addRuleRegex: regexp.MustCompile(botConfig.AddRuleRegex),
		config: botConfig,
	}
}

func (sb *SwearBot) LoadSwears() {
	file, err := os.Open(sb.dictFileName)
	if err != nil {
		log.Fatalf("Error opening swear dictionary file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := normalizeWord(scanner.Text())
		sb.dict.AddEntry(word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from swear dictionary file: %v", err)
	}
}

func (sb *SwearBot) ParseMessage(message string) string {
	rules := sb.addRuleRegex.FindAllStringSubmatch(message, 1)
	if rules != nil {
		rule := rules[0][1]
		err := sb.addRule(rule)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf(sb.config.OnAddRuleResponse, rule)
	}
	return sb.parseSwears(message)
}

func (sb *SwearBot) addRule(rule string) error {
	file, fileReadErr := os.OpenFile(sb.dictFileName, os.O_RDWR | os.O_APPEND, 0666)
	if fileReadErr != nil {
		log.Printf("Add rule: Cannot open swear dictionary file: %v", fileReadErr)
		return errors.New(sb.config.OnAddRuleFileReadErr)
	}
	defer file.Close()

	normRule := normalizeWord(rule)

	confilctErr := sb.dict.AddEntry(normRule)
	if confilctErr != nil {
		log.Printf("Add rule: %s", confilctErr.Desc)
		return errors.New(sb.config.OnAddRuleConflictErr)
	}

	_, saveErr := file.WriteString(fmt.Sprintf("%s\n", normRule))
	if saveErr != nil {
		log.Printf("Add rule: Cannot write string '%s' to swear dictionary file: %v", normRule, saveErr)
		return errors.New(sb.config.OnAddRuleSaveErr)
	}

	return nil
}

func (sb *SwearBot) parseSwears(message string) string {
	swears := sb.findSwears(message)
	if len(swears) > 0 {
		swearsLine := fmt.Sprintf("*%s*", strings.Join(swears, "*, *"))
		response := fmt.Sprintf(sb.config.OnSwearsFoundResponse, swearsLine)
		return response
	}
	return ""
}

func (sb *SwearBot) findSwears(message string) []string {
	swears := make([]string, 0)
	words := strings.Fields(message)
	for _, word := range words {
		word = normalizeWord(word)
		success, _ := sb.dict.Match(word)
		if success {
			swears = append(swears, word)
		}
	}
	return swears
}

func normalizeWord(word string) string {
	word = strings.Trim(word, " \n\r")
	word = strings.ToLower(word)
	return word
}


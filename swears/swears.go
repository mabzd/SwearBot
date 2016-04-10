package swears

import (
	"os"
	"log"
	"bufio"
	"errors"
	"fmt"
	"strings"
	"../dictmatch"
)

type Swears struct {
	dict *dictmatch.Dict
	dictFileName string
	config SwearsConfig
}

type SwearsConfig struct {
	OnAddRuleFileReadErr string
	OnAddRuleConflictErr string
	OnAddRuleSaveErr string
	OnIvalidWildcardErr string
}

func NewSwears(dictFileName string, config SwearsConfig) *Swears {
	return &Swears {
		dict: dictmatch.NewDict(),
		dictFileName: dictFileName,
		config: config,
	}
}

func (sw *Swears) LoadSwears() {
	file, err := os.Open(sw.dictFileName)
	if err != nil {
		log.Fatalf("Error opening swear dictionary file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := normalizeWord(scanner.Text())
		sw.dict.AddEntry(word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from swear dictionary file: %v", err)
	}
}

func (sw *Swears) AddRule(rule string) error {
	file, fileReadErr := os.OpenFile(sw.dictFileName, os.O_RDWR | os.O_APPEND, 0666)
	if fileReadErr != nil {
		log.Printf("Add rule: Cannot open swear dictionary file: %v", fileReadErr)
		return errors.New(sw.config.OnAddRuleFileReadErr)
	}
	defer file.Close()

	normRule := normalizeWord(rule)

	confilctErr := sw.dict.AddEntry(normRule)
	if confilctErr != nil {
		log.Printf("Add rule: %s", confilctErr.Desc)
		if (confilctErr.ErrType == dictmatch.InvalidWildardPlacementErr) {
			return errors.New(sw.config.OnIvalidWildcardErr)
		}
		return errors.New(sw.config.OnAddRuleConflictErr)
	}

	_, saveErr := file.WriteString(fmt.Sprintf("%s\n", normRule))
	if saveErr != nil {
		log.Printf("Add rule: Cannot write string '%s' to swear dictionary file: %v", normRule, saveErr)
		return errors.New(sw.config.OnAddRuleSaveErr)
	}

	return nil
}

func (sw *Swears) FindSwears(message string) []string {
	swears := make([]string, 0)
	words := strings.Fields(message)
	for _, word := range words {
		word = normalizeWord(word)
		success, _ := sw.dict.Match(word)
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
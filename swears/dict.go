package swears

import (
	"../dictmatch"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	DictFileReadErr    = 21
	InvalidWildcardErr = 22
	AddRuleConflictErr = 23
	AddRuleSaveErr     = 24
)

func (sw *Swears) AddRule(rule string) int {
	file, fileReadErr := os.OpenFile(sw.config.DictFileName, os.O_RDWR|os.O_APPEND, 0666)
	if fileReadErr != nil {
		log.Printf("Add rule: Cannot open swear dictionary file: %v", fileReadErr)
		return DictFileReadErr
	}
	defer file.Close()

	normRule := normalizeWord(rule)

	confilctErr := sw.dict.AddEntry(normRule)
	if confilctErr != nil {
		log.Printf("Add rule: %s", confilctErr.Desc)
		if confilctErr.ErrType == dictmatch.InvalidWildardPlacementErr {
			return InvalidWildcardErr
		}

		return AddRuleConflictErr
	}

	_, saveErr := file.WriteString(fmt.Sprintf("%s\n", normRule))
	if saveErr != nil {
		log.Printf("Add rule: Cannot write string '%s' to swear dictionary file: %v", normRule, saveErr)
		return AddRuleSaveErr
	}

	return Success
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

func (sw *Swears) LoadSwears() int {
	file, err := os.Open(sw.config.DictFileName)
	if err != nil {
		log.Printf("Load swears: Error opening swear dictionary file: %v", err)
		return DictFileReadErr
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := normalizeWord(scanner.Text())
		sw.dict.AddEntry(word)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Load swears: Error reading from swear dictionary file: %v", err)
		return DictFileReadErr
	}

	return Success
}

func normalizeWord(word string) string {
	word = strings.Trim(word, " \n\r")
	word = strings.ToLower(word)
	return word
}

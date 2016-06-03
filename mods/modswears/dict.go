package modswears

import (
	"../../dictmatch"
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

func (mod *ModSwears) AddRule(rule string) int {
	file, fileReadErr := os.OpenFile(mod.dictFileName, os.O_RDWR|os.O_APPEND, 0666)
	if fileReadErr != nil {
		log.Printf("ModSwears: cannot open swear dictionary file: %v\n", fileReadErr)
		return DictFileReadErr
	}
	defer file.Close()

	normRule := normalizeWord(rule)

	confilctErr := mod.dict.AddEntry(normRule)
	if confilctErr != nil {
		log.Printf("ModSwears: add rule: %s\n", confilctErr.Desc)
		if confilctErr.ErrType == dictmatch.InvalidWildardPlacementErr {
			return InvalidWildcardErr
		}

		return AddRuleConflictErr
	}

	_, saveErr := file.WriteString(fmt.Sprintf("%s\n", normRule))
	if saveErr != nil {
		log.Printf("ModSwears: cannot write string '%s' to swear dictionary file: %v\n", normRule, saveErr)
		return AddRuleSaveErr
	}

	return Success
}

func (mod *ModSwears) FindSwears(message string) []string {
	swears := make([]string, 0)
	words := strings.Fields(message)
	for _, word := range words {
		word = normalizeWord(word)
		success, _ := mod.dict.Match(word)
		if success {
			swears = append(swears, word)
		}
	}

	return swears
}

func (mod *ModSwears) LoadSwears() int {
	file, err := os.Open(mod.dictFileName)
	if err != nil {
		log.Printf("ModSwears: Error opening swear dictionary file: %v\n", err)
		return DictFileReadErr
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := normalizeWord(scanner.Text())
		mod.dict.AddEntry(word)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("ModSwears: Error reading from swear dictionary file: %v\n", err)
		return DictFileReadErr
	}

	return Success
}

func normalizeWord(word string) string {
	word = strings.Trim(word, " \n\r")
	word = strings.ToLower(word)
	return word
}

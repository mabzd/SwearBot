package swearbot

import (
	"log"
	"os"
	"bufio"
	"strings"
	"../dictmatch"
)

type SwearBot struct {
	dict *dictmatch.Dict
	dictFileName string
}

func NewSwearBot(fileName string) *SwearBot {
	return &SwearBot {
		dict: dictmatch.NewDict(),
		dictFileName: fileName,
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

func (sb *SwearBot) FindSwears(message string) []string {
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


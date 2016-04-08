package swearbot

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"regexp"
	"../dictmatch"
)

var AddSwearRegex *regexp.Regexp = regexp.MustCompile("/^\\s*add rule: ([a-z*]+)\\s*$/i")

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

func (sb *SwearBot) ParseMessage(message string) string {
	swears := sb.findSwears(message)
	if len(swears) > 0 {
		response := fmt.Sprintf("Following swears found: *%s*", strings.Join(swears, "*, *"))
		return response
	}
	return ""
}

func (sb *SwearBot) addSwear(swear string) {
	sb.addSwears([]string { swear })
}

func (sb *SwearBot) addSwears(swears []string) {
	file, err := os.OpenFile(sb.dictFileName, os.O_RDWR | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening swear dictionary file: %v", err)
	}
	defer file.Close()

	for _, swear := range swears {
		file.WriteString(fmt.Sprintf("%s\n", normalizeWord(swear)))
	}
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

func addSwearToFile(swear string) {

}

func normalizeWord(word string) string {
	word = strings.Trim(word, " \n\r")
	word = strings.ToLower(word)
	return word
}


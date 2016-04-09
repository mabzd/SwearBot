package swearbot

import (
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"testing"
)

func TestSwears(t *testing.T) {
	config := createBotConfig()
	tmpfile := createTempDict()
	defer os.Remove(tmpfile.Name())

	sb := createBot(tmpfile, config)
	assertResponse(t, sb, "Test Aa\tA Abc\rB\nAbba", "Swears: *a*, *abba*")
}

func TestAddRule(t *testing.T) {
	config := createBotConfig()
	tmpfile := createTempDict()
	defer os.Remove(tmpfile.Name())

	sb := createBot(tmpfile, config)
	assertResponse(t, sb, " Add rule:   XXX* ", "Rule: XXX*")
	assertResponse(t, sb, "Test ABBA\tXx\nXxxxx", "Swears: *abba*, *xxxxx*")
}

func TestAddRuleFileReadErr(t *testing.T) {
	config := createBotConfig()
	tmpfile := createTempDict()
	defer os.Remove(tmpfile.Name())

	sb := createBot(tmpfile, config)
	os.Remove(tmpfile.Name())
	assertResponse(t, sb, "add rule: r1", "FileReadErr")
}

func TestAddRuleConflictErr(t *testing.T) {
	config := createBotConfig()
	tmpfile := createTempDict()
	defer os.Remove(tmpfile.Name())

	sb := createBot(tmpfile, config)
	assertResponse(t, sb, "add rule: ab*", "ConflictErr")
}

func TestAddRuleInvalidWildcardErr(t *testing.T) {
	config := createBotConfig()
	tmpfile := createTempDict()
	defer os.Remove(tmpfile.Name())

	sb := createBot(tmpfile, config)
	assertResponse(t, sb, "add rule: a*b", "InvalidWildcard")
}

func createBotConfig() BotConfig {
	return BotConfig {
		AddRuleRegex: "(?i)^\\s*add rule:\\s*([a-z0-9*]+)\\s*$",
		OnSwearsFoundResponse: "Swears: %s",
		OnAddRuleResponse: "Rule: %s",
		OnAddRuleFileReadErr: "FileReadErr",
		OnAddRuleConflictErr: "ConflictErr",
		OnAddRuleSaveErr: "SaveErr",
		OnIvalidWildcardErr: "InvalidWildcard",
	}
}

func createTempDict() *os.File {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}
	loadTempDict(tmpfile)
	tmpfile.Close()
	return tmpfile
}

func loadTempDict(file *os.File) {
	rules := []string { "a", "abcd", "abb*" }
	for _, rule := range rules {
		file.WriteString(fmt.Sprintf("%s\n", rule))
	}
}

func createBot(file *os.File, config BotConfig) *SwearBot {
	sb := NewSwearBot(file.Name(), config)
	sb.LoadSwears()
	return sb
}

func assertResponse(t *testing.T, sb *SwearBot, message string, response string) {
	r := sb.ParseMessage(message)
	if response != r {
		t.Fatalf("Expected add rule response '%s' but got '%s'", response, r)
	}
}
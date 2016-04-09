package swearbot

import (
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"testing"
)

func TestSwearBot(t *testing.T) {
	config := createBotConfig()
	tmpfile := createTempDict()
	defer os.Remove(tmpfile.Name())

	sb := NewSwearBot(tmpfile.Name(), config)
	sb.LoadSwears()
	response := sb.ParseMessage("test aa\ta abc\rb\nabba")

	expectedResponse := "Swears: *a*, *abba*"
	if response != expectedResponse {
		t.Fatalf("Expected response '%s' but got '%s'", expectedResponse, response)
	}
}

func createBotConfig() BotConfig {
	return BotConfig {
		AddRuleRegex: "/^\\s*add rule: ([a-z*]+)\\s*$/i",
		OnSwearsFoundResponse: "Swears: %s",
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
package bot

import (
	"testing"
)

func TestMentions(t *testing.T) {
	compileMentionRegex("bot.*")

	m := "<@bot.*>"
	assertIsMention(t, m)
	assertRemoveMentions(t, m, "")

	m = "<@bot.*>: this is message"
	assertIsMention(t, m)
	assertRemoveMentions(t, m, " this is message")

	m = "<@bot.*>: this <@bot.*> is message <@bot.*>:"
	assertIsMention(t, m)
	assertRemoveMentions(t, m, " this  is message ")

	assertIsNotMention(t, "bot.*")
	assertIsNotMention(t, "bot.*:")
	assertIsNotMention(t, "@bot.")
	assertIsNotMention(t, "@bot.:")
	assertIsNotMention(t, "<@bott>")
	assertIsNotMention(t, "<@bot>")
}

func assertIsMention(t *testing.T, message string) {
	if !isMention(message) {
		t.Fatalf("Mention not found in message '%s'", message)
	}
}

func assertIsNotMention(t *testing.T, message string) {
	if isMention(message) {
		t.Fatalf("Expected no mentions in message '%s'", message)
	}
}

func assertRemoveMentions(t *testing.T, message string, expected string) {
	actual := removeMentions(message)
	if actual != expected {
		t.Fatalf("Expected processed message %#v, got %#v", expected, actual)
	}
}

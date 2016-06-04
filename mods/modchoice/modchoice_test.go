package modchoice

import (
	"reflect"
	"testing"
)

func TestGetOptions(t *testing.T) {
	assertGetOptions(t, "o1 or o2", []string{"o1", "o2"})
	assertGetOptions(t, "o1 | o2", []string{"o1", "o2"})
	assertGetOptions(t, "o1 or o2?", []string{"o1", "o2"})
	assertGetOptions(t, "o1 | o2?", []string{"o1", "o2"})
	assertGetOptions(t, "o1 or o2? ", []string{"o1", "o2"})
	assertGetOptions(t, "o1 | o2? ", []string{"o1", "o2"})
	assertGetOptions(t, "o1 or o2 ? ", []string{"o1", "o2"})
	assertGetOptions(t, "o1 | o2 ? ", []string{"o1", "o2"})
	assertGetOptions(t, "o1oro2", []string{"o1oro2"})
	assertGetOptions(t, "o1|o2", []string{"o1|o2"})
	assertGetOptions(t, "  o1  or   o2  ", []string{"o1", "o2"})
	assertGetOptions(t, "  o1  |   o2  ", []string{"o1", "o2"})
	assertGetOptions(t, "o1 or o2 | o3", []string{"o1", "o2", "o3"})
	assertGetOptions(t, "o1 | o2 or o3", []string{"o1", "o2", "o3"})
	assertGetOptions(t, "o1 or o2 | o3 or o4", []string{"o1", "o2", "o3", "o4"})
	assertGetOptions(t, "o1 | o2 or o3 | o4", []string{"o1", "o2", "o3", "o4"})
	assertGetOptions(t, "o1 or ", []string{"o1"})
	assertGetOptions(t, "o1 | ", []string{"o1"})
	assertGetOptions(t, "o1 or", []string{"o1"})
	assertGetOptions(t, "o1 |", []string{"o1"})
	assertGetOptions(t, "or o1", []string{"o1"})
	assertGetOptions(t, "| o1", []string{"o1"})
	assertGetOptions(t, " or o1", []string{"o1"})
	assertGetOptions(t, " | o1", []string{"o1"})
	assertGetOptions(t, "o1 or or o2", []string{"o1", "o2"})
	assertGetOptions(t, "o1 | | o2", []string{"o1", "o2"})
	assertGetOptions(t, "o1 or | o2", []string{"o1", "o2"})
	assertGetOptions(t, "o1 | or o2", []string{"o1", "o2"})
	assertGetOptions(t, "o1 oror o2", []string{"o1 oror o2"})
	assertGetOptions(t, "o1 || o2", []string{"o1 || o2"})
	assertGetOptions(t, "or", []string{})
	assertGetOptions(t, "|", []string{})
}

func TestModChoiceResponse(t *testing.T) {
	mod := createModChoice()
	assertProcessMention(t, mod, "a or b", []string{"r1 a", "r1 b", "r2"})
	assertProcessMention(t, mod, "a | b", []string{"r1 a", "r1 b", "r2"})
}

func TestModChoiceNoResponse(t *testing.T) {
	mod := createModChoice()
	assertProcessMention(t, mod, "test test", []string{""})
	assertProcessMention(t, mod, "testortest", []string{""})
	assertProcessMention(t, mod, "test|test", []string{""})
	assertProcessMention(t, mod, "a |  ", []string{""})
	assertProcessMention(t, mod, "or", []string{""})
	assertProcessMention(t, mod, "or|", []string{""})
}

func TestModNullChoiceResponse(t *testing.T) {
	mod := createModChoice()
	mod.config.NullChoiceProbability = 1.0
	assertProcessMention(t, mod, "a or b", []string{"n1", "n2"})
}

func TestModAllResponses(t *testing.T) {
	mod := createModChoice()
	mod.config.NullChoiceProbability = 0.5
	assertProcessMention(t, mod, "a or b", []string{"r1 a", "r1 b", "r2", "n1", "n2"})
	assertProcessMention(t, mod, "a | b", []string{"r1 a", "r1 b", "r2", "n1", "n2"})
}

func createModChoice() *ModChoice {
	mod := NewModChoice()
	mod.config = &ModChoiceConfig{
		OrKeywords:            []string{"or", "|"},
		ChoiceResponseFormat:  []string{"r1 {option}", "r2"},
		NullChoiceProbability: 0.0,
		NullChoiceResponses:   []string{"n1", "n2"},
	}
	return mod
}

func assertGetOptions(t *testing.T, message string, expected []string) {
	orKeys := []string{"or", "|"}
	actual := getOptions(message, orKeys)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(
			"Expected options %#v from message '%s', got %#v",
			expected,
			message,
			actual)
	}
}

func assertProcessMention(
	t *testing.T,
	mod *ModChoice,
	message string,
	possibleExpected []string) {

	actual := mod.ProcessMention(message, "u1", "c1")
	for _, expected := range possibleExpected {
		if actual == expected {
			return
		}
	}
	t.Errorf(
		"Expected one response from %#v when processing message '%s', got '%s'",
		possibleExpected,
		message,
		actual)
}

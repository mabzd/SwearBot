package modmention

import (
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestValidateConfig(t *testing.T) {
	mod := NewModMention()
	mod.config = &ModMentionConfig{
		Reactions: []*Reaction{
			&Reaction{
				Weight:    0,
				Responses: []string{"r1", "r2"},
			},
			&Reaction{
				Weight:    1001,
				Responses: []string{"r3", "r4"},
			},
			&Reaction{
				Weight:    10,
				Responses: []string{"r5"},
			},
			&Reaction{
				Weight:    10,
				Responses: []string{},
			},
		},
	}
	expected := &ModMentionConfig{
		Reactions: []*Reaction{
			&Reaction{
				Weight:    2,
				Responses: []string{"r1", "r2"},
			},
			&Reaction{
				Weight:    2002,
				Responses: []string{"r3", "r4"},
			},
			&Reaction{
				Weight:    2012,
				Responses: []string{"r5"},
			},
		},
	}
	assertValidConfig(t, mod)
	assertMaxWeight(t, mod, 2012)
	assertConfigEq(t, mod, expected)
}

func TestValidateEmptyConfig(t *testing.T) {
	mod := NewModMention()
	mod.config = &ModMentionConfig{
		Reactions: []*Reaction{
			&Reaction{
				Weight:    1,
				Responses: []string{},
			},
		},
	}
	assertInvalidConfig(t, mod)
}

func TestGetReactions(t *testing.T) {
	mod := NewModMention()
	mod.config = &ModMentionConfig{
		Reactions: []*Reaction{
			&Reaction{
				Weight:    1,
				Responses: []string{"r1", "r2"},
			},
			&Reaction{
				Weight:    100,
				Responses: []string{"r3", "r4"},
			},
			&Reaction{
				Weight:    10,
				Responses: []string{"r5"},
			},
		},
	}
	assertValidConfig(t, mod)
	assertMaxWeight(t, mod, 212)
	assertGetReaction(t, mod, 0, []string{"r1", "r2"})
	assertGetReaction(t, mod, 1, []string{"r1", "r2"})
	assertGetReaction(t, mod, 2, []string{"r3", "r4"})
	assertGetReaction(t, mod, 201, []string{"r3", "r4"})
	assertGetReaction(t, mod, 202, []string{"r5"})
	assertGetReaction(t, mod, 211, []string{"r5"})
	assertGetReaction(t, mod, 212, []string{""})
}

func assertValidConfig(t *testing.T, mod *ModMention) {
	if !mod.validateConfig() {
		t.Fatal("Expected valid config")
	}
}

func assertInvalidConfig(t *testing.T, mod *ModMention) {
	if mod.validateConfig() {
		t.Fatal("Expected invalid config")
	}
}

func assertConfigEq(t *testing.T, mod *ModMention, expected *ModMentionConfig) {
	if !reflect.DeepEqual(mod.config, expected) {
		t.Fatal("Config deep equal failed")
	}
}

func assertMaxWeight(t *testing.T, mod *ModMention, expected int) {
	if mod.maxWeight != expected {
		t.Fatalf("Expected maxWeight to be %d, got %d", expected, mod.maxWeight)
	}
}

func assertGetReaction(t *testing.T, mod *ModMention, d int, possibleExpected []string) {
	actual := mod.getReaction(d)
	for _, expected := range possibleExpected {
		if actual == expected {
			return
		}
	}
	t.Fatalf(
		"Expected one reaction from %#v for distribution %d, got '%s'",
		possibleExpected,
		d,
		actual)
}

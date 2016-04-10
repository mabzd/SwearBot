package swearbot

import (
	"os"
	"io/ioutil"
	"testing"
	"reflect"
)

func TestAddSwears(t *testing.T) {
	tmpFilePath := createTmpFilePath(t)
	defer os.Remove(tmpFilePath)

	sb := createStatsBot(tmpFilePath)
	assertAddSwears(t, sb, 1, 2016, "user1", 3)
	assertAddSwears(t, sb, 1, 2016, "user1", 2)

	expected := []*User {
		&User {
			Name: "user1",
			SwearCount: 5,
		},
	}

	assertMonthlyRank(t, sb, 1, 2016, expected)
}

func TestRankOrder(t *testing.T) {
	tmpFilePath := createTmpFilePath(t)
	defer os.Remove(tmpFilePath)

	sb := createStatsBot(tmpFilePath)
	assertAddSwears(t, sb, 1, 2016, "user1", 3)
	assertAddSwears(t, sb, 1, 2016, "user2", 4)
	assertAddSwears(t, sb, 1, 2016, "user1", 2)
	assertAddSwears(t, sb, 1, 2016, "user3", 6)
	assertAddSwears(t, sb, 2, 2016, "user1", 10)

	expected := []*User {
		&User {
			Name: "user3",
			SwearCount: 6,
		},
		&User {
			Name: "user1",
			SwearCount: 5,
		},
		&User {
			Name: "user2",
			SwearCount: 4,
		},
	}

	assertMonthlyRank(t, sb, 1, 2016, expected)
}

func TestUnknownMonth(t *testing.T) {
	tmpFilePath := createTmpFilePath(t)
	defer os.Remove(tmpFilePath)

	sb := createStatsBot(tmpFilePath)
	assertAddSwears(t, sb, 1, 2016, "user1", 1)

	assertMonthlyRank(t, sb, 2, 2016, []*User {})
}

func createTmpFilePath(t *testing.T) string {
	tmpfile, err := ioutil.TempFile("", "stats")
	if err != nil {
		t.Fatalf("Cannot create tmp file: %s", err)
	}
	path := tmpfile.Name()
	tmpfile.Close()
	os.Remove(tmpfile.Name())
	return path
}

func createStatsBot(tmpFilePath string) *SwearBot {
	return NewSwearBot("", tmpFilePath, BotConfig {})
}

func assertAddSwears(t *testing.T, sb *SwearBot, m int, y int, u string, n int) {
	err := sb.AddSwears(m, y, u, n)
	if err != nil {
		t.Fatalf("Expected no error when adding swears but got %s", err)
	}
}

func assertMonthlyRank(t *testing.T, sb *SwearBot, m int, y int, expected []*User) {
	users, err := sb.GetMonthlyRank(m, y)
	if err != nil {
		t.Fatalf("Expected no error when getting monthly rank but got %s", err)
	}
	if !reflect.DeepEqual(users, expected) {
		t.Fatal("Monthly rank deep equal failed")
	}
}
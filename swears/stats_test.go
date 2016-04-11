package swears

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestAddSwears(t *testing.T) {
	tmpFilePath := createTmpFilePath(t)
	defer os.Remove(tmpFilePath)

	sw := createStats(tmpFilePath)
	assertAddSwearCount(t, sw, 1, 2016, "user1", 3)
	assertAddSwearCount(t, sw, 1, 2016, "user1", 2)

	expected := []*UserStats{
		&UserStats{
			UserId:     "user1",
			SwearCount: 5,
		},
	}

	assertMonthlyRank(t, sw, 1, 2016, expected)
}

func TestRankOrder(t *testing.T) {
	tmpFilePath := createTmpFilePath(t)
	defer os.Remove(tmpFilePath)

	sw := createStats(tmpFilePath)
	assertAddSwearCount(t, sw, 1, 2016, "user1", 3)
	assertAddSwearCount(t, sw, 1, 2016, "user2", 4)
	assertAddSwearCount(t, sw, 1, 2016, "user1", 2)
	assertAddSwearCount(t, sw, 1, 2016, "user3", 6)
	assertAddSwearCount(t, sw, 2, 2016, "user1", 10)

	expected := []*UserStats{
		&UserStats{
			UserId:     "user3",
			SwearCount: 6,
		},
		&UserStats{
			UserId:     "user1",
			SwearCount: 5,
		},
		&UserStats{
			UserId:     "user2",
			SwearCount: 4,
		},
	}

	assertMonthlyRank(t, sw, 1, 2016, expected)
}

func TestUnknownMonth(t *testing.T) {
	tmpFilePath := createTmpFilePath(t)
	defer os.Remove(tmpFilePath)

	sw := createStats(tmpFilePath)
	assertAddSwearCount(t, sw, 1, 2016, "user1", 1)

	assertMonthlyRank(t, sw, 2, 2016, []*UserStats{})
}

func createTmpFilePath(t *testing.T) string {
	tmpfile, err := ioutil.TempFile("", "Swears")
	if err != nil {
		t.Fatalf("Cannot create tmp file: %s", err)
	}
	path := tmpfile.Name()
	tmpfile.Close()
	os.Remove(tmpfile.Name())
	return path
}

func createStats(tmpFilePath string) *Swears {
	config := SwearsConfig{
		StatsFileName: tmpFilePath,
	}
	return NewSwears(nil, config)
}

func assertAddSwearCount(t *testing.T, sw *Swears, m int, y int, u string, n int) {
	err := sw.AddSwearCount(m, y, u, n)
	if err != nil {
		t.Fatalf("Expected no error when adding swears but got %s", err)
	}
}

func assertMonthlyRank(t *testing.T, sw *Swears, m int, y int, expected []*UserStats) {
	users, err := sw.GetMonthlyRank(m, y)
	if err != nil {
		t.Fatalf("Expected no error when getting monthly rank but got %s", err)
	}
	if !reflect.DeepEqual(users, expected) {
		t.Fatal("Monthly rank deep equal failed")
	}
}

package swears

import (
	"../utils"
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
	tmpFilePath := createTmpStatsPath(t)
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
	tmpFilePath := createTmpStatsPath(t)
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
	tmpFilePath := createTmpStatsPath(t)
	defer os.Remove(tmpFilePath)

	sw := createStats(tmpFilePath)
	assertAddSwearCount(t, sw, 1, 2016, "user1", 1)

	assertMonthlyRank(t, sw, 2, 2016, []*UserStats{})
}

func TestTotalRank(t *testing.T) {
	tmpFilePath := createTmpStatsPath(t)
	defer os.Remove(tmpFilePath)

	sw := createStats(tmpFilePath)
	assertAddSwearCount(t, sw, 1, 2016, "user1", 1)
	assertAddSwearCount(t, sw, 1, 2016, "user2", 1)
	assertAddSwearCount(t, sw, 2, 2016, "user1", 2)
	assertAddSwearCount(t, sw, 3, 2016, "user1", 1)
	assertAddSwearCount(t, sw, 3, 2016, "user2", 4)
	assertAddSwearCount(t, sw, 3, 2016, "user3", 3)

	expected := []*UserStats{
		&UserStats{
			UserId:     "user2",
			SwearCount: 5,
		},
		&UserStats{
			UserId:     "user1",
			SwearCount: 4,
		},
		&UserStats{
			UserId:     "user3",
			SwearCount: 3,
		},
	}

	assertTotalRank(t, sw, expected)
}

func TestEmptyTotalRank(t *testing.T) {
	tmpFilePath := createTmpStatsPath(t)
	defer os.Remove(tmpFilePath)

	sw := createStats(tmpFilePath)
	assertTotalRank(t, sw, []*UserStats{})
}

func createTmpStatsPath(t *testing.T) string {
	fileName := utils.CreateTmpFileName("Stats")
	if fileName == "" {
		t.Fatal("Cannot create temp stats file path")
	}
	return fileName
}

func createStats(tmpFilePath string) *Swears {
	config := SwearsConfig{
		StatsFileName: tmpFilePath,
	}
	return NewSwears(nil, config)
}

func assertAddSwearCount(t *testing.T, sw *Swears, m int, y int, u string, n int) {
	err := sw.AddSwearCount(m, y, u, n)
	if err != Success {
		t.Fatalf("Expected no error when adding swears but got %v", err)
	}
}

func assertMonthlyRank(t *testing.T, sw *Swears, m int, y int, expected []*UserStats) {
	actual, err := sw.GetMonthlyRank(m, y)
	if err != Success {
		t.Fatalf("Expected no error when getting monthly rank but got %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("Monthly rank deep equal failed")
	}
}

func assertTotalRank(t *testing.T, sw *Swears, expected []*UserStats) {
	actual, err := sw.GetTotalRank()
	if err != Success {
		t.Fatalf("Expected no error when getting total rank but got %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("Total rank deep equal failed")
	}
}

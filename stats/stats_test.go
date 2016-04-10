package stats

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestAddSwears(t *testing.T) {
	tmpFilePath := createTmpFilePath(t)
	defer os.Remove(tmpFilePath)

	st := createStats(tmpFilePath)
	assertAddSwearCount(t, st, 1, 2016, "user1", 3)
	assertAddSwearCount(t, st, 1, 2016, "user1", 2)

	expected := []*User{
		&User{
			Name:       "user1",
			SwearCount: 5,
		},
	}

	assertMonthlyRank(t, st, 1, 2016, expected)
}

func TestRankOrder(t *testing.T) {
	tmpFilePath := createTmpFilePath(t)
	defer os.Remove(tmpFilePath)

	st := createStats(tmpFilePath)
	assertAddSwearCount(t, st, 1, 2016, "user1", 3)
	assertAddSwearCount(t, st, 1, 2016, "user2", 4)
	assertAddSwearCount(t, st, 1, 2016, "user1", 2)
	assertAddSwearCount(t, st, 1, 2016, "user3", 6)
	assertAddSwearCount(t, st, 2, 2016, "user1", 10)

	expected := []*User{
		&User{
			Name:       "user3",
			SwearCount: 6,
		},
		&User{
			Name:       "user1",
			SwearCount: 5,
		},
		&User{
			Name:       "user2",
			SwearCount: 4,
		},
	}

	assertMonthlyRank(t, st, 1, 2016, expected)
}

func TestUnknownMonth(t *testing.T) {
	tmpFilePath := createTmpFilePath(t)
	defer os.Remove(tmpFilePath)

	st := createStats(tmpFilePath)
	assertAddSwearCount(t, st, 1, 2016, "user1", 1)

	assertMonthlyRank(t, st, 2, 2016, []*User{})
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

func createStats(tmpFilePath string) *Stats {
	return NewStats(tmpFilePath, StatsConfig{})
}

func assertAddSwearCount(t *testing.T, st *Stats, m int, y int, u string, n int) {
	err := st.AddSwearCount(m, y, u, n)
	if err != nil {
		t.Fatalf("Expected no error when adding swears but got %s", err)
	}
}

func assertMonthlyRank(t *testing.T, st *Stats, m int, y int, expected []*User) {
	users, err := st.GetMonthlyRank(m, y)
	if err != nil {
		t.Fatalf("Expected no error when getting monthly rank but got %s", err)
	}
	if !reflect.DeepEqual(users, expected) {
		t.Fatal("Monthly rank deep equal failed")
	}
}

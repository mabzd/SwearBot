package utils

import (
	"testing"
	"time"
)

func TestLastDayOfPrevMonth(t *testing.T) {
	assertLastDayOfPrevMonth(t, getDate(2016, 6, 15), getDate(2016, 5, 31))
	assertLastDayOfPrevMonth(t, getDate(2016, 6, 1), getDate(2016, 5, 31))
	assertLastDayOfPrevMonth(t, getDate(2016, 3, 31), getDate(2016, 2, 29))
	assertLastDayOfPrevMonth(t, getDate(2016, 1, 3), getDate(2015, 12, 31))
}

func getDate(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location())
}

func assertLastDayOfPrevMonth(t *testing.T, date time.Time, expected time.Time) {
	actual := LastDayOfPrevMonth(date)
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

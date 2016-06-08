package utils

import (
	"time"
)

var TimeClock Clock = RealClock{}

type Clock interface {
	Now() time.Time
}

type RealClock struct{}
type MockClock struct {
	CurrentTime time.Time
}

func (c RealClock) Now() time.Time {
	return time.Now()
}

func (c MockClock) Now() time.Time {
	return c.CurrentTime
}

func NewLocalDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func NewLocalDateTime(year int, month time.Month, day int, h int, m int, s int) time.Time {
	return time.Date(year, month, day, h, m, s, 0, time.Local)
}

func LastDayOfPrevMonth(date time.Time) time.Time {
	month := date.Month()
	year := date.Year()
	return time.Date(year, month, 1, 0, 0, 0, 0, date.Location()).AddDate(0, 0, -1)
}

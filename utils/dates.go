package utils

import (
	"time"
)

func LastDayOfPrevMonth(date time.Time) time.Time {
	month := date.Month()
	year := date.Year()
	return time.Date(year, month, 1, 0, 0, 0, 0, date.Location()).AddDate(0, 0, -1)
}

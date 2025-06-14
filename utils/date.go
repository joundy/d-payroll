package utils

import "time"

var TimeNow = time.Now

func GetStartOfDay() time.Time {
	now := TimeNow()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

func GetEndOfDay() time.Time {
	now := TimeNow()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
}

func IsWeekend() bool {
	t := TimeNow()

	switch t.Weekday() {
	case time.Saturday, time.Sunday:
		return true
	default:
		return false
	}
}

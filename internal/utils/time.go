package utils

import "time"

func NextWeekRange(loc *time.Location) (time.Time, time.Time) {
	start := time.Now().In(loc)
	todayStart := StartOfDay(start)
	end := todayStart.AddDate(0, 0, 7)

	return start.UTC(), end.UTC()
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

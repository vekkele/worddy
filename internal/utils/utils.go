package utils

import (
	"time"
)

func HumanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func CalculateNextReview(hoursToNext int32) time.Time {
	currentTime := time.Now()
	dur := time.Hour * time.Duration(hoursToNext)

	return currentTime.Add(dur).Truncate(time.Hour)
}

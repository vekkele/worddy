package utils

import (
	"time"
)

func CalculateNextReview(hoursToNext int32) time.Time {
	currentTime := time.Now()
	dur := time.Hour * time.Duration(hoursToNext)

	return currentTime.Add(dur).Truncate(time.Hour)
}

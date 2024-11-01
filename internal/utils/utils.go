package utils

import (
	"math"
	"time"
)

func HumanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

type Stage string

const (
	Apprentice  = "apprentice"
	Guru        = "guru"
	Master      = "master"
	Enlightened = "enlightened"
	Burned      = "burned"
)

var stageColors = map[Stage]string{
	Apprentice:  "#db2777",
	Guru:        "#9333ea",
	Master:      "#2563eb",
	Enlightened: "#0284c7",
	Burned:      "#57534e",
}

func GetStageColor(stage Stage) string {
	return stageColors[stage]
}

func GetStageFromLevel(level int32) Stage {
	switch level {
	case 9:
		return Burned
	case 8:
		return Enlightened
	case 7:
		return Master
	case 5, 6:
		return Guru
	default:
		return Apprentice
	}
}

func CalculateNextReview(hoursToNext int32) time.Time {
	currentTime := time.Now()
	dur := time.Hour * time.Duration(hoursToNext)

	return currentTime.Add(dur).Truncate(time.Hour)
}

func CalculateNextStage(curLevel int32, wrongAnswers int32) int32 {
	if wrongAnswers == 0 {
		return curLevel + 1
	}

	incorrectAdjustmentCount := int32(math.Ceil(float64(wrongAnswers) / 2.0))
	var srsPenaltyFactor int32 = 1
	if curLevel >= 5 {
		srsPenaltyFactor = 2
	}

	nextLevel := curLevel - incorrectAdjustmentCount*srsPenaltyFactor
	if nextLevel < 1 {
		return 1
	}

	return nextLevel
}

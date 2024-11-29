package domain

import (
	"math"
	"strings"
	"time"
)

type Stage string

const (
	Apprentice  Stage = "apprentice"
	Guru        Stage = "guru"
	Master      Stage = "master"
	Enlightened Stage = "enlightened"
	Burned      Stage = "burned"
)

type Word struct {
	ID           int64
	Word         string
	Translations []string
	NextReview   time.Time
	StageLevel   int32
}

var stageColors = map[Stage]string{
	Apprentice:  "#db2777",
	Guru:        "#9333ea",
	Master:      "#2563eb",
	Enlightened: "#0284c7",
	Burned:      "#57534e",
}

func (w *Word) StageColor() string {
	return stageColors[w.StageName()]
}

func (w *Word) StageName() Stage {
	switch w.StageLevel {
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

func (w *Word) NextStage(wrongAnswers int32) int32 {
	if wrongAnswers == 0 {
		return w.StageLevel + 1
	}

	incorrectAdjustmentCount := int32(math.Ceil(float64(wrongAnswers) / 2.0))
	var srsPenaltyFactor int32 = 1
	if w.StageLevel >= 5 {
		srsPenaltyFactor = 2
	}

	nextLevel := w.StageLevel - incorrectAdjustmentCount*srsPenaltyFactor
	if nextLevel < 1 {
		return 1
	}

	return nextLevel
}

func (w *Word) CheckTranslation(guess string) bool {
	clearedGuess := strings.ToLower(strings.TrimSpace(guess))

	for _, translation := range w.Translations {
		clearedTranslation := strings.ToLower(translation)

		if clearedTranslation == clearedGuess {
			return true
		}
	}

	return false
}

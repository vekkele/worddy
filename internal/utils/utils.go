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

type Stage string

const (
	Apprentice  = "apprentice"
	Guru        = "guru"
	Master      = "master"
	Enlightened = "enlightened"
	Burned      = "burned"
)

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

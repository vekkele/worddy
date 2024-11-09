package domain

import (
	"time"

	"github.com/vekkele/worddy/internal/utils"
)

type Word struct {
	ID           int64
	Word         string
	Translations []string
	NextReview   time.Time
	StageLevel   int32
	StageName    utils.Stage
}

package store

import (
	"context"

	"github.com/vekkele/worddy/internal/domain"
)

type ReviewWordsStore interface {
	InitReview(ctx context.Context, userID int64) error
	GetNextReviewWord(ctx context.Context, userID int64) (domain.ReviewWord, error)
	CheckWord(ctx context.Context, userID int64, wordID int64, guess string) (bool, []string, error)
}

package store

import (
	"context"

	"github.com/vekkele/worddy/internal/domain"
)

type WordStore interface {
	Insert(ctx context.Context, userID int64, word string, translations []string) error
	GetAll(ctx context.Context, userID int64) ([]domain.Word, error)
	GetReview(ctx context.Context, userID int64) ([]domain.Word, error)
	UpdateWordStage(ctx context.Context, id, userID int64, wrongAnswers int32) error
}

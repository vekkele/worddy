package store

import (
	"context"
	"time"

	"github.com/vekkele/worddy/internal/domain"
)

type UserStore interface {
	Insert(ctx context.Context, email, passwordHash string) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Exists(ctx context.Context, id int64) (bool, error)
}

type WordStore interface {
	Insert(ctx context.Context, userID int64, word string, translations []string) error
	GetAll(ctx context.Context, userID int64) ([]domain.Word, error)
	GetByID(ctx context.Context, userID int64, wordID int64) (domain.Word, error)
	UpdateWordLevel(ctx context.Context, userID, wordID int64, level int32) error

	InitReview(ctx context.Context, userID int64) error
	GetNextReviewWord(ctx context.Context, userID int64) (domain.ReviewWord, error)
	CommitWrongAnswer(ctx context.Context, userID, wordID int64) error
	GetWrongAnswers(ctx context.Context, userID, wordID int64) (int32, error)
	DeleteReviewWord(ctx context.Context, userID, wordID int64) error

	GetReviewsCountInRange(ctx context.Context, userID int64, start time.Time, end time.Time) ([]domain.ReviewsAtTime, error)

	GetUserReviewCount(ctx context.Context, userID int64) (int, error)
}

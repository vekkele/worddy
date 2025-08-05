package mocks

import (
	"context"
	"errors"
	"time"

	"github.com/vekkele/worddy/internal/domain"
	"github.com/vekkele/worddy/internal/store"
)

type wordStore struct{}

func NewWordStore() store.WordStore {
	return &wordStore{}
}

func (m *wordStore) Insert(ctx context.Context, userID int64, word string, translations []string) error {
	return nil
}

func (m *wordStore) GetAll(ctx context.Context, userID int64) ([]domain.Word, error) {
	return []domain.Word{
		{ID: 1, Word: "Word1", Translations: []string{"translation1, translation2"}},
		{ID: 2, Word: "Word2", Translations: []string{"translation3, translation4"}},
	}, nil
}

func (m *wordStore) GetReview(ctx context.Context, userID int64) ([]domain.Word, error) {
	return []domain.Word{
		{ID: 1, Word: "Word1", Translations: []string{"translation1, translation2"}},
		{ID: 2, Word: "Word2", Translations: []string{"translation3, translation4"}},
	}, nil
}

func (m *wordStore) CommitWrongAnswer(ctx context.Context, userID int64, wordID int64) error {
	return nil
}

func (m *wordStore) DeleteReviewWord(ctx context.Context, userID int64, wordID int64) error {
	return nil
}

func (m *wordStore) GetByID(ctx context.Context, userID int64, wordID int64) (domain.Word, error) {
	if wordID == 1 && userID == 1 {
		return domain.Word{
			ID:           1,
			Word:         "Word1",
			Translations: []string{"translation1, translation2"},
		}, nil
	}

	return domain.Word{}, errors.New("service mock: no word found")
}

func (m *wordStore) GetNextReviewWord(ctx context.Context, userID int64) (domain.ReviewWord, error) {
	return domain.ReviewWord{
		Word: domain.Word{
			ID:           1,
			Word:         "Word1",
			Translations: []string{"translation1, translation2"},
		},
		WrongAnswers: 0,
	}, nil
}

func (m *wordStore) GetWrongAnswers(ctx context.Context, userID int64, wordID int64) (int32, error) {
	return 0, nil
}

func (m *wordStore) InitReview(ctx context.Context, userID int64) error {
	return nil
}

func (m *wordStore) UpdateWordLevel(ctx context.Context, userID int64, wordID int64, level int32) error {
	return nil
}

func (m *wordStore) GetReviewsCountInRange(ctx context.Context, userID int64, start time.Time, end time.Time) ([]domain.ReviewsAtTime, error) {
	return nil, nil
}

func (m *wordStore) GetUserReviewCount(ctx context.Context, userID int64) (int, error) {
	return 0, nil
}

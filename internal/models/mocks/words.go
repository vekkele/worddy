package mocks

import (
	"context"

	"github.com/vekkele/worddy/internal/models"
)

type WordModel struct{}

func (m *WordModel) Insert(ctx context.Context, userID int64, word string, translations []string) error {
	return nil
}

func (m *WordModel) GetAll(ctx context.Context, userID int64) ([]models.Word, error) {
	return []models.Word{
		{ID: 1, Word: "Word1", Translations: "translation1, translation2"},
		{ID: 2, Word: "Word2", Translations: "translation3, translation4"},
	}, nil
}

func (m *WordModel) GetReview(ctx context.Context, userID int64) ([]models.Word, error) {
	return []models.Word{
		{ID: 1, Word: "Word1", Translations: "translation1, translation2"},
		{ID: 2, Word: "Word2", Translations: "translation3, translation4"},
	}, nil
}

func (m *WordModel) UpdateWordStage(ctx context.Context, id, userID int64, wrongAnswers int32) error {
	return nil
}

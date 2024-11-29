package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/vekkele/worddy/internal/domain"
	"github.com/vekkele/worddy/internal/store/postgres/db"
)

func (m *wordStore) InitReview(ctx context.Context, userID int64) error {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := m.db.WithTx(tx)

	words, err := qtx.GetUserReviewWords(ctx, userID)
	if err != nil {
		return err
	}

	for _, word := range words {
		err := qtx.AddReviewWord(ctx, db.AddReviewWordParams{
			UserID: userID,
			WordID: word.ID,
		})

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *wordStore) GetNextReviewWord(ctx context.Context, userID int64) (domain.ReviewWord, error) {
	row, err := s.db.GetNextReviewWord(ctx, userID)
	if err != nil {
		return domain.ReviewWord{}, err
	}

	wordRow, err := s.db.GetWordByID(ctx, db.GetWordByIDParams{UserID: userID, ID: row.WordID})
	if err != nil {
		return domain.ReviewWord{}, err
	}

	return getReviewWordFromDB(wordRow, row), nil
}

func (s *wordStore) DeleteReviewWord(ctx context.Context, userID, wordID int64) error {
	return s.db.DeleteReviewWord(ctx, db.DeleteReviewWordParams{WordID: wordID, UserID: userID})
}

func (s *wordStore) CommitWrongAnswer(ctx context.Context, userID, wordID int64) error {
	return s.db.CommitWrongReviewAnswer(ctx, db.CommitWrongReviewAnswerParams{WordID: wordID, UserID: userID})
}

func (s *wordStore) GetWrongAnswers(ctx context.Context, userID, wordID int64) (int32, error) {
	return s.db.GetReviewWrongAnswers(ctx, db.GetReviewWrongAnswersParams{WordID: wordID, UserID: userID})
}

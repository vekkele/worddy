package models

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vekkele/worddy/internal/models/db"
	"github.com/vekkele/worddy/internal/utils"
)

type ReviewWordsModel interface {
	InitReview(ctx context.Context, userID int64) error
	GetNextReviewWord(ctx context.Context, userID int64) (ReviewWord, error)
	CheckWord(ctx context.Context, userID int64, wordID int64, guess string) (bool, []string, error)
}

type reviewWordsModel struct {
	db   *db.Queries
	pool *pgxpool.Pool
}

func NewReviewWordsModel(pool *pgxpool.Pool) ReviewWordsModel {
	queries := db.New(pool)
	return &reviewWordsModel{db: queries, pool: pool}
}

func (m *reviewWordsModel) InitReview(ctx context.Context, userID int64) error {
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

type ReviewWord struct {
	Word
	WrongAnswers int32
}

func (m *reviewWordsModel) GetNextReviewWord(ctx context.Context, userID int64) (ReviewWord, error) {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return ReviewWord{}, err
	}
	defer tx.Rollback(ctx)

	qtx := m.db.WithTx(tx)

	row, err := qtx.GetNextReviewWord(ctx, userID)
	if err != nil {
		return ReviewWord{}, err
	}

	wordRow, err := qtx.GetWordByID(ctx, db.GetWordByIDParams{UserID: userID, ID: row.WordID})
	if err != nil {
		return ReviewWord{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return ReviewWord{}, err
	}

	return ReviewWord{
		Word:         getWordFromDBRow(db.GetUserWordsRow(wordRow)),
		WrongAnswers: row.WrongAnswers,
	}, nil
}

func (m *reviewWordsModel) CheckWord(ctx context.Context, userID int64, wordID int64, guess string) (bool, []string, error) {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, nil, err
	}
	defer tx.Rollback(ctx)

	qtx := m.db.WithTx(tx)

	row, err := qtx.GetWordByID(ctx, db.GetWordByIDParams{
		UserID: userID, ID: wordID,
	})
	if err != nil {
		return false, nil, err
	}

	word := getWordFromDBRow(db.GetUserWordsRow(row))

	right := utils.CheckTranslation(word.Translations, guess)

	if right {
		err = updateReviewedWord(ctx, qtx, word.StageLevel, wordID, userID)
		if err != nil {
			return false, nil, err
		}
	} else {

		err = qtx.CommitWrongReviewAnswer(ctx, db.CommitWrongReviewAnswerParams{WordID: wordID, UserID: userID})
		if err != nil {
			return false, nil, err
		}
	}

	err = tx.Commit(ctx)

	return right, word.Translations, err
}

func updateReviewedWord(ctx context.Context, q *db.Queries, currentStage int32, wordID int64, userID int64) error {
	reviewWord, err := q.GetReviewWordByID(ctx, db.GetReviewWordByIDParams{WordID: wordID, UserID: userID})
	if err != nil {
		return err
	}

	nextLevel := utils.CalculateNextStage(currentStage, reviewWord.WrongAnswers)

	stage, err := q.GetStageByLevel(ctx, nextLevel)
	if err != nil {
		return err
	}

	err = q.UpdateWordStage(ctx, db.UpdateWordStageParams{
		ID:      wordID,
		UserID:  userID,
		StageID: stage.ID,
		NextReview: pgtype.Timestamptz{
			Time:  utils.CalculateNextReview(stage.HoursToNext),
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	return q.DeleteReviewWord(ctx, db.DeleteReviewWordParams{WordID: wordID, UserID: userID})
}

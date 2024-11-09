package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vekkele/worddy/internal/domain"
	"github.com/vekkele/worddy/internal/store"
	"github.com/vekkele/worddy/internal/store/postgres/db"
	"github.com/vekkele/worddy/internal/utils"
)

type wordStore struct {
	db   *db.Queries
	pool *pgxpool.Pool
}

func NewWordStore(pool *pgxpool.Pool) store.WordStore {
	db := db.New(pool)
	return &wordStore{
		db:   db,
		pool: pool,
	}
}

func (m *wordStore) Insert(ctx context.Context, userID int64, word string, translations []string) error {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := m.db.WithTx(tx)

	stage, err := qtx.GetStageByLevel(ctx, 1)
	if err != nil {
		return err
	}

	var nextReviewTimestampz pgtype.Timestamptz

	err = nextReviewTimestampz.Scan(utils.CalculateNextReview(stage.HoursToNext))
	if err != nil {
		return err
	}

	id, err := qtx.AddWord(ctx, db.AddWordParams{
		Word:       word,
		NextReview: nextReviewTimestampz,
		StageID:    stage.ID,
		UserID:     userID,
	})
	if err != nil {
		return err
	}

	for _, t := range translations {
		err = qtx.AddTranslation(ctx, db.AddTranslationParams{WordID: id, Translation: t})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (m *wordStore) GetAll(ctx context.Context, userID int64) ([]domain.Word, error) {
	rows, err := m.db.GetUserWords(ctx, userID)
	if err != nil {
		return nil, err
	}

	var words []domain.Word
	for _, w := range rows {
		words = append(words, getWordFromDBRow(w))
	}

	return words, nil
}

func (m *wordStore) GetReview(ctx context.Context, userID int64) ([]domain.Word, error) {
	rows, err := m.db.GetUserReviewWords(ctx, userID)
	if err != nil {
		return nil, err
	}

	var words []domain.Word
	for _, w := range rows {
		words = append(words, getWordFromDBRow(db.GetUserWordsRow(w)))
	}

	return words, nil
}

func getWordFromDBRow(row db.GetUserWordsRow) domain.Word {
	return domain.Word{
		ID:           row.ID,
		Word:         row.Word,
		Translations: row.Translations,
		NextReview:   row.NextReview.Time,
		StageLevel:   row.Level,
		StageName:    utils.GetStageFromLevel(row.Level),
	}
}

func (m *wordStore) UpdateWordStage(ctx context.Context, id, userID int64, wrongAnswers int32) error {
	word, err := m.db.GetWordByID(ctx, db.GetWordByIDParams{ID: id, UserID: userID})
	if err != nil {
		return err
	}

	nextLevel := utils.CalculateNextStage(word.Level, wrongAnswers)

	stage, err := m.db.GetStageByLevel(ctx, nextLevel)
	if err != nil {
		return err
	}

	return m.db.UpdateWordStage(ctx, db.UpdateWordStageParams{
		ID:      id,
		UserID:  userID,
		StageID: stage.ID,
		NextReview: pgtype.Timestamptz{
			Time:  utils.CalculateNextReview(stage.HoursToNext),
			Valid: true,
		},
	})
}

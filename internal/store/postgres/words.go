package postgres

import (
	"context"
	"time"

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
		Word: word,
		NextReview: pgtype.Timestamptz{
			Time:  utils.CalculateNextReview(stage.HoursToNext),
			Valid: true,
		},
		StageID: stage.ID,
		UserID:  userID,
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
		words = append(words, getWordFromDBRow(db.GetWordByIDRow(w)))
	}

	return words, nil
}

func (s *wordStore) GetByID(ctx context.Context, userID int64, wordID int64) (domain.Word, error) {
	row, err := s.db.GetWordByID(ctx, db.GetWordByIDParams{ID: wordID, UserID: userID})
	if err != nil {
		return domain.Word{}, err
	}

	return getWordFromDBRow(row), err
}

func (m *wordStore) UpdateWordLevel(ctx context.Context, userID, wordID int64, level int32) error {
	stage, err := m.db.GetStageByLevel(ctx, level)
	if err != nil {
		return err
	}

	return m.db.UpdateWordStage(ctx, db.UpdateWordStageParams{
		ID:      wordID,
		UserID:  userID,
		StageID: stage.ID,
		NextReview: pgtype.Timestamptz{
			Time:  utils.CalculateNextReview(stage.HoursToNext),
			Valid: true,
		},
	})
}

func (s *wordStore) GetReviewsCountInRange(ctx context.Context, userID int64, start time.Time, end time.Time) ([]domain.ReviewsAtTime, error) {
	rows, err := s.db.GetUserReviewsCountInRange(ctx, db.GetUserReviewsCountInRangeParams{
		UserID:       userID,
		NextReview:   pgtype.Timestamptz{Time: start, Valid: true},
		NextReview_2: pgtype.Timestamptz{Time: end, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	if len(rows) < 1 {
		return []domain.ReviewsAtTime{}, nil
	}

	var reviews []domain.ReviewsAtTime
	for _, row := range rows {
		reviews = append(reviews, domain.ReviewsAtTime{
			Count: int(row.Count),
			Time:  row.NextReview.Time,
		})
	}

	return reviews, nil
}

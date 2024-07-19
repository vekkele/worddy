package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vekkele/worddy/internal/models/db"
)

type WordModel interface {
	Insert(ctx context.Context, userID int64, word string, translations []string) error
}

type wordModel struct {
	db   *db.Queries
	pool *pgxpool.Pool
}

func NewWordModel(pool *pgxpool.Pool) WordModel {
	db := db.New(pool)
	return &wordModel{
		db:   db,
		pool: pool,
	}
}

func (m *wordModel) Insert(ctx context.Context, userID int64, word string, translations []string) error {
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

	err = nextReviewTimestampz.Scan(m.calculateNextReview(int(stage.HoursToNext)))
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

func (m *wordModel) calculateNextReview(hoursToNext int) time.Time {
	currentTime := time.Now()
	dur := time.Hour * time.Duration(hoursToNext)

	return currentTime.Add(dur)
}

package service

import (
	"context"
	"time"

	"github.com/vekkele/worddy/internal/domain"
	"github.com/vekkele/worddy/internal/store"
	"github.com/vekkele/worddy/internal/utils"
)

type WordService interface {
	Add(ctx context.Context, userID int64, word string, translations []string) error
	GetAll(ctx context.Context, userID int64) ([]domain.Word, error)
	InitReview(ctx context.Context, userID int64) (domain.ReviewWord, error)
	CheckWord(ctx context.Context, userID, wordID int64, guess string) (bool, []string, error)
	GetReviewForecast(ctx context.Context, userID int64, timeZone string) ([]domain.DayForecast, error)
	GetReviewCount(ctx context.Context, userID int64) (int, error)
}

type wordService struct {
	store store.WordStore
}

func NewWordService(store store.WordStore) WordService {
	return &wordService{store: store}
}

func (s *wordService) Add(ctx context.Context, userID int64, word string, translations []string) error {
	return s.store.Insert(ctx, userID, word, translations)
}

func (s *wordService) GetAll(ctx context.Context, userID int64) ([]domain.Word, error) {
	return s.store.GetAll(ctx, userID)
}

func (s *wordService) GetReviewCount(ctx context.Context, userID int64) (int, error) {
	return s.store.GetUserReviewCount(ctx, userID)
}

func (s *wordService) InitReview(ctx context.Context, userID int64) (domain.ReviewWord, error) {
	err := s.store.InitReview(ctx, userID)
	if err != nil {
		return domain.ReviewWord{}, err
	}

	return s.store.GetNextReviewWord(ctx, userID)
}

func (s *wordService) CheckWord(ctx context.Context, userID, wordID int64, guess string) (bool, []string, error) {
	word, err := s.store.GetByID(ctx, userID, wordID)
	if err != nil {
		return false, nil, err
	}

	right := word.CheckTranslation(guess)

	if right {
		err = s.updateWordStage(ctx, userID, word)
	} else {
		err = s.store.CommitWrongAnswer(ctx, userID, wordID)
	}

	if err != nil {
		return false, nil, err
	}

	return right, word.Translations, nil
}

func (s *wordService) updateWordStage(ctx context.Context, userID int64, word domain.Word) error {
	wrongAnswers, err := s.store.GetWrongAnswers(ctx, userID, word.ID)
	if err != nil {
		return err
	}

	nextLevel := word.NextStage(wrongAnswers)

	err = s.store.UpdateWordLevel(ctx, userID, word.ID, nextLevel)
	if err != nil {
		return err
	}

	return s.store.DeleteReviewWord(ctx, userID, word.ID)
}

func (s *wordService) GetReviewForecast(ctx context.Context, userID int64, timeZone string) ([]domain.DayForecast, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return nil, err
	}

	start, end := utils.NextWeekRange(loc)

	reviews, err := s.store.GetReviewsCountInRange(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}

	return domain.NewForecast(reviews, loc), nil
}

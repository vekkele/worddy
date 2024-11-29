package postgres

import (
	"github.com/vekkele/worddy/internal/domain"
	"github.com/vekkele/worddy/internal/store/postgres/db"
)

func getWordFromDBRow(row db.GetWordByIDRow) domain.Word {
	return domain.Word{
		ID:           row.ID,
		Word:         row.Word,
		Translations: row.Translations,
		NextReview:   row.NextReview.Time,
		StageLevel:   row.Level,
	}
}

func getReviewWordFromDB(wordRow db.GetWordByIDRow, wrongAnswersRow db.GetNextReviewWordRow) domain.ReviewWord {
	return domain.ReviewWord{
		Word:         getWordFromDBRow(wordRow),
		WrongAnswers: wrongAnswersRow.WrongAnswers,
	}
}

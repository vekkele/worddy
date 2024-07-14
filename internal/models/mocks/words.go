package mocks

import "context"

type WordModel struct{}

func (m *WordModel) Insert(ctx context.Context, userID int64, word string, translations []string) error {
	return nil
}

package mocks

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vekkele/worddy/internal/models"
	"github.com/vekkele/worddy/internal/store/postgres/db"
)

type UserModel struct{}

var User = db.User{
	ID:    1,
	Email: "john@doe.test",
	CreatedAt: pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	},
}

const UserPassword = "pass123"

func (m *UserModel) Insert(ctx context.Context, email, password string) error {
	if email == User.Email {
		return models.ErrDuplicateEmail
	}

	return nil
}

func (m *UserModel) Authenticate(ctx context.Context, email, password string) (int64, error) {
	if email == User.Email && password == UserPassword {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(ctx context.Context, id int64) (bool, error) {
	if id == User.ID {
		return true, nil
	}

	return false, nil
}

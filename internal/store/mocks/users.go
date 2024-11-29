package mocks

import (
	"context"
	"log"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vekkele/worddy/internal/domain"
	"github.com/vekkele/worddy/internal/store"
	"github.com/vekkele/worddy/internal/store/postgres/db"
)

type userStore struct{}

func NewUserStore() store.UserStore {
	return &userStore{}
}

const UserPassword = "pass123"

func createMockUser() db.User {
	hash, err := argon2id.CreateHash(UserPassword, argon2id.DefaultParams)
	if err != nil {
		log.Fatalf("mock: failed to hash password")
	}

	return db.User{
		ID:           1,
		Email:        "john@doe.test",
		PasswordHash: []byte(hash),
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

}

var User = createMockUser()

func (m *userStore) Insert(ctx context.Context, email, password string) error {
	if email == User.Email {
		return domain.ErrDuplicateEmail
	}

	return nil
}

func (m *userStore) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == User.Email {
		return &domain.User{
			ID: User.ID, PasswordHash: string(User.PasswordHash),
		}, nil
	}

	return nil, domain.ErrNoUserFound
}

func (m *userStore) Exists(ctx context.Context, id int64) (bool, error) {
	if id == User.ID {
		return true, nil
	}

	return false, nil
}

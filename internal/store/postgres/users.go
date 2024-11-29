package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vekkele/worddy/internal/domain"
	"github.com/vekkele/worddy/internal/store"
	"github.com/vekkele/worddy/internal/store/postgres/db"
)

type userStore struct {
	db   *db.Queries
	pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) store.UserStore {
	queries := db.New(pool)
	return &userStore{db: queries, pool: pool}
}

func (m *userStore) Insert(ctx context.Context, email, passwordHash string) error {
	err := m.db.CreateUser(
		ctx, db.CreateUserParams{
			Email:        email,
			PasswordHash: []byte(passwordHash),
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
				return domain.ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (m *userStore) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := m.db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNoUserFound
		}
		return nil, err
	}

	return &domain.User{
		ID:           row.ID,
		PasswordHash: string(row.PasswordHash),
	}, nil
}

func (m *userStore) Exists(ctx context.Context, id int64) (bool, error) {
	return m.db.UserExists(ctx, id)
}

package models

import (
	"context"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vekkele/worddy/internal/models/db"
)

type UserModel struct {
	db   db.Queries
	pool *pgxpool.Pool
}

func NewUserModel(pool *pgxpool.Pool) UserModel {
	queries := db.New(pool)
	return UserModel{db: *queries, pool: pool}
}

func (m *UserModel) Insert(ctx context.Context, email, password string) error {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	err = m.db.CreateUser(
		ctx, db.CreateUserParams{
			Email:        email,
			PasswordHash: []byte(hash),
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
				return ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (m *UserModel) Authenticate(ctx context.Context, email, password string) (int64, error) {
	row, err := m.db.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, string(row.PasswordHash))
	if err != nil {
		return 0, err
	}

	if !match {
		return 0, ErrInvalidCredentials
	}

	return row.ID, nil
}

func (m *UserModel) Exists(ctx context.Context, id int64) (bool, error) {
	return m.db.Exists(ctx, id)
}

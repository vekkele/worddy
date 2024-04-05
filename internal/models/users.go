package models

import (
	"context"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5"
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

func (m *UserModel) Insert(ctx context.Context, email, password string) (db.User, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return db.User{}, err
	}

	user, err := m.db.CreateUser(
		ctx, db.CreateUserParams{
			Email:        email,
			PasswordHash: []byte(hash),
		},
	)
	if err != nil {
		//TODO: Return specific error if email is already taken
		return db.User{}, err
	}

	return user, nil
}

func (m *UserModel) Authenticate(ctx context.Context, email, password string) (int64, error) {
	row, getErr := m.db.GetByEmail(ctx, email)

	match, err := argon2id.ComparePasswordAndHash(password, string(row.PasswordHash))
	if err != nil {
		return 0, err
	}

	if getErr != nil {
		if errors.Is(getErr, pgx.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, getErr
	}

	if !match {
		return 0, ErrInvalidCredentials
	}

	return row.ID, nil
}

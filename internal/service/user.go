package service

import (
	"context"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/vekkele/worddy/internal/domain"
	"github.com/vekkele/worddy/internal/store"
)

type UserService interface {
	Register(ctx context.Context, email, password string) error
	Authenticate(ctx context.Context, email, password string) (int64, error)
	Exists(ctx context.Context, id int64) (bool, error)
}

type userService struct {
	store store.UserStore
}

func NewUserService(store store.UserStore) UserService {
	return &userService{store: store}
}

func (s *userService) Register(ctx context.Context, email, password string) error {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	return s.store.Insert(ctx, email, hash)
}

func (s *userService) Authenticate(ctx context.Context, email, password string) (int64, error) {
	user, err := s.store.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNoUserFound) {
			return 0, domain.ErrInvalidCredentials
		}
		return 0, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return 0, err
	}

	if !match {
		return 0, domain.ErrInvalidCredentials
	}

	return user.ID, nil
}

func (s *userService) Exists(ctx context.Context, id int64) (bool, error) {
	return s.store.Exists(ctx, id)
}

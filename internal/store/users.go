package store

import (
	"context"
)

type UserStore interface {
	Insert(ctx context.Context, email, password string) error
	Authenticate(ctx context.Context, email, password string) (int64, error)
	Exists(ctx context.Context, id int64) (bool, error)
}

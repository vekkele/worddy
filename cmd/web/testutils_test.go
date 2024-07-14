package main

import (
	"io"
	"log/slog"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/vekkele/worddy/internal/models/mocks"
)

func newTestApplication() *application {
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour

	return &application{
		users:          &mocks.UserModel{},
		words:          &mocks.WordModel{},
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

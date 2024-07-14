package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
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

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

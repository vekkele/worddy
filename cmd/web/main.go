package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/vekkele/worddy/internal/config"
	"github.com/vekkele/worddy/internal/service"
	"github.com/vekkele/worddy/internal/store/postgres"
)

type application struct {
	users          service.UserService
	words          service.WordService
	formDecoder    *form.Decoder
	logger         *slog.Logger
	sessionManager *scs.SessionManager
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	config, err := config.New()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	pool, err := postgres.OpenDB(config.DB.DSN)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	userStore := postgres.NewUserStore(pool)
	wordStore := postgres.NewWordStore(pool)

	userService := service.NewUserService(userStore)
	wordService := service.NewWordService(wordStore)

	app := application{
		users:          userService,
		words:          wordService,
		formDecoder:    form.NewDecoder(),
		logger:         logger,
		sessionManager: sessionManager,
	}

	err = http.ListenAndServe(":"+config.Port, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

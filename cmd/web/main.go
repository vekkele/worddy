package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	"github.com/vekkele/worddy/internal/config"
	"github.com/vekkele/worddy/internal/models"
	"github.com/vekkele/worddy/internal/postgres"
)

type application struct {
	users       models.UserModel
	formDecoder *form.Decoder
	logger      *slog.Logger
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

	app := application{
		users:       models.NewUserModel(pool),
		formDecoder: form.NewDecoder(),
		logger:      logger,
	}

	err = http.ListenAndServe(":"+config.Port, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

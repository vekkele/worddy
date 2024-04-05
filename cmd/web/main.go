package main

import (
	"log"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/vekkele/worddy/internal/config"
	"github.com/vekkele/worddy/internal/models"
	"github.com/vekkele/worddy/internal/postgres"
)

type application struct {
	userModel   models.UserModel
	formDecoder *form.Decoder
}

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := postgres.OpenDB(config.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	userModel := models.NewUserModel(pool)

	app := application{
		userModel:   userModel,
		formDecoder: form.NewDecoder(),
	}

	err = http.ListenAndServe(":"+config.Port, app.routes())
	log.Fatal(err)
}

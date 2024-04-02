package main

import (
	"log"
	"net/http"

	"github.com/vekkele/worddy/internal/config"
	"github.com/vekkele/worddy/internal/postgres"
)

type application struct{}

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	_, err = postgres.OpenDB(config.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	app := application{}

	err = http.ListenAndServe(":"+config.Port, app.routes())
	log.Fatal(err)
}

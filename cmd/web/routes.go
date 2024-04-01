package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vekkele/worddy/ui"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	fs := http.FileServer(http.FS(ui.Files))
	r.Handle("/static/*", fs)

	r.Get("/", app.home)

	r.Get("/", app.home)
	r.Get("/user/signup", app.signup)
	r.Get("/user/login", app.login)

	return r
}

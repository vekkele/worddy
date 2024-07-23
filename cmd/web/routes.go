package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vekkele/worddy/ui"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.recoverPanic, app.logRequest, secureHeaders)

	fs := http.FileServer(http.FS(ui.Files))
	r.Handle("/static/*", fs)

	r.Group(func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

		r.Get("/", app.home)
		r.Get("/user/signup", app.signup)
		r.Post("/user/signup", app.signupPost)
		r.Get("/user/login", app.login)
		r.Post("/user/login", app.loginPost)

		r.Group(func(r chi.Router) {
			r.Use(app.requireAuthentication)

			r.Post("/user/logout", app.logoutPost)
			r.Get("/dashboard", app.dashboard)
			r.Get("/word/add", app.wordAdd)
			r.Post("/word/add", app.wordAddPost)
			r.Get("/review", app.review)
			r.Put("/review", app.reviewPost)
		})
	})

	return r
}

package main

import (
	"net/http"

	"github.com/vekkele/worddy/ui/view/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	render(w, r, pages.Home("Home"))
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	render(w, r, pages.Signup())
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	render(w, r, pages.Login())
}

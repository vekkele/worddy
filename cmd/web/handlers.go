package main

import (
	"fmt"
	"net/http"

	"github.com/vekkele/worddy/ui/view/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pages.Home("Home"))
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pages.Signup())
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pages.Login())
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	var form pages.SignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//TODO: Validate form
	user, err := app.users.Insert(r.Context(), form.Email, form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	user.PasswordHash = nil

	fmt.Fprintf(w, "User created: %#v\n", user)
}

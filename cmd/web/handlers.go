package main

import (
	"errors"
	"net/http"

	"github.com/vekkele/worddy/internal/models"
	"github.com/vekkele/worddy/internal/validator"
	"github.com/vekkele/worddy/ui/view/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pages.Home("Home"))
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pages.Signup(pages.SignupForm{}))
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

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Invalid email format")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 6), "password", "This field must be at least 6 characters")
	form.CheckField(form.Password == form.PasswordConfirm, "password-confirm", "Passwords do not match")

	if !form.Valid() {
		app.render(w, r, pages.Signup(form))
		return
	}

	if err := app.users.Insert(r.Context(), form.Email, form.Password); err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			app.render(w, r, pages.Signup(form))
			return
		}

		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

package main

import (
	"errors"
	"net/http"

	"github.com/vekkele/worddy/internal/models"
	"github.com/vekkele/worddy/internal/validator"
	"github.com/vekkele/worddy/ui/view/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	app.render(w, r, pages.Home(r, "Home"))
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pages.Signup(r, pages.SignupForm{}))
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pages.Login(r, pages.LoginForm{}))
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
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, pages.Signup(r, form))
		return
	}

	if err := app.users.Insert(r.Context(), form.Email, form.Password); err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(w, r, pages.Signup(r, form))
			return
		}

		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	var form pages.LoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Invalid email format")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	userID, err := app.users.Authenticate(r.Context(), form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Invalid email or password")
			app.render(w, r, pages.Login(r, form))
			return
		}

		app.serverError(w, r, err)
		return
	}

	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", userID)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {
	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pages.Dashboard(r))
}

func (app *application) wordAdd(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pages.WordAdd(r, pages.WordAddForm{}))
}

func (app *application) wordAddPost(w http.ResponseWriter, r *http.Request) {
	var form pages.WordAddForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	translations := splitTranslations(form.Translations)

	form.CheckField(validator.NotBlank(form.Word), "word", "This field cannot be blank")
	form.CheckField(len(translations) > 0, "translations", "At least one non-empty translation must be provided")

	if !form.Valid() {
		app.render(w, r, pages.WordAdd(r, form))
		return
	}

	userID := app.authenticatedUserID(r)
	if userID == 0 {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	err = app.words.Insert(r.Context(), userID, form.Word, translations)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

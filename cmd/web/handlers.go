package main

import (
	"errors"
	"net/http"

	"github.com/vekkele/worddy/internal/domain"
	"github.com/vekkele/worddy/internal/validator"
	"github.com/vekkele/worddy/ui/view/pages"
	"github.com/vekkele/worddy/ui/view/partials"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	app.render(w, r, pages.Home(r, "Home"))
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	app.render(w, r, pages.SignupPage(r, pages.SignupFormData{}))
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	app.render(w, r, pages.LoginPage(r, pages.LoginFormData{}))
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	var form pages.SignupFormData

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
		app.render(w, r, pages.SignupForm(r, form))
		return
	}

	if err := app.users.Register(r.Context(), form.Email, form.Password); err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(w, r, pages.SignupForm(r, form))
			return
		}

		app.serverError(w, r, err)
		return
	}

	w.Header().Set("HX-Location", "/user/login")
	w.WriteHeader(http.StatusSeeOther)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	var form pages.LoginFormData

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Invalid email format")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, pages.LoginForm(r, form))
		return
	}

	userID, err := app.users.Authenticate(r.Context(), form.Email, form.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			form.AddNonFieldError("Invalid email or password")
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(w, r, pages.LoginForm(r, form))
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

	w.Header().Set("HX-Location", "/dashboard")
	w.WriteHeader(http.StatusSeeOther)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {
	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	userID := app.authenticatedUserID(r)
	words, err := app.words.GetAll(r.Context(), userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, pages.Dashboard(r, words))
}

func (app *application) wordAdd(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pages.WordAdd(r, pages.WordAddFormData{}))
}

func (app *application) wordAddPost(w http.ResponseWriter, r *http.Request) {
	var form pages.WordAddFormData

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	translations := splitTranslations(form.Translations)

	form.CheckField(validator.NotBlank(form.Word), "word", "This field cannot be blank")
	form.CheckField(len(translations) > 0, "translations", "At least one non-empty translation must be provided")

	if !form.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, pages.WordAddForm(r, form))
		return
	}

	userID := app.authenticatedUserID(r)

	err = app.words.Add(r.Context(), userID, form.Word, translations)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("HX-Location", "/dashboard")
	w.WriteHeader(http.StatusSeeOther)
}

func (app *application) review(w http.ResponseWriter, r *http.Request) {
	userID := app.authenticatedUserID(r)
	nextWord, err := app.words.InitReview(r.Context(), userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, pages.Review(r, nextWord))
}

func (app *application) checkWord(w http.ResponseWriter, r *http.Request) {
	var form partials.CheckWordForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	userID := app.authenticatedUserID(r)

	correct, translations, err := app.words.CheckWord(r.Context(), userID, form.WordID, form.Guess)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, partials.WordCheckResult(partials.WordCheckResultProps{
		CheckFormProps: partials.CheckFormProps{
			WordID:       form.WordID,
			Checked:      true,
			CorrectGuess: correct,
			GuessValue:   form.Guess,
		},
		Translations: translations,
	}))
}

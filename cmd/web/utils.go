package main

import (
	"errors"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-playground/form/v4"
	"github.com/vekkele/worddy/ui/view/partials"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderErr *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderErr) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) renderError(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Set("HX-Retarget", "#toast-section")
	w.Header().Set("HX-Reswap", "innerHTML")
	app.render(w, r, partials.ToastError(message))
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()
	trace := string(debug.Stack())

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)

	app.renderError(w, r, "Something went wrong. Please try again")
}

func (app *application) authenticatedUserID(r *http.Request) int64 {
	return app.sessionManager.GetInt64(r.Context(), "authenticatedUserID")
}

func (app *application) saveSession(r *http.Request, userID int64) error {
	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		return err
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", userID)
	return nil
}

func (app *application) clearSession(r *http.Request) error {
	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		return err
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func splitTranslations(raw string) []string {
	translations := []string{}
	for _, tr := range strings.Split(raw, ",") {
		trimmed := strings.TrimSpace(tr)

		if trimmed != "" {
			translations = append(translations, trimmed)
		}
	}

	return translations
}

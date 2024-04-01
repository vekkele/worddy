package main

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

func render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

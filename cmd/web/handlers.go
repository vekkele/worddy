package main

import (
	"github.com/labstack/echo/v4"
	"github.com/vekkele/worddy/ui/view"
)

func (app *application) home(c echo.Context) error {
	return render(c, view.Index("Worddy"))
}

package main

import (
	"github.com/labstack/echo/v4"
	"github.com/vekkele/worddy/ui/view/pages"
)

func (app *application) home(c echo.Context) error {
	return render(c, pages.Home("Worddy"))
}

func (app *application) signup(c echo.Context) error {
	return render(c, pages.Signup())
}

func (app *application) login(c echo.Context) error {
	return render(c, pages.Login())
}

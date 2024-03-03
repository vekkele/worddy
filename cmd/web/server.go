package main

import "github.com/labstack/echo/v4"

func (app *application) start() {
	e := echo.New()

	e.GET("/", app.home)

	e.Logger.Fatal(e.Start(app.addr))
}

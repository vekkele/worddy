package main

import "github.com/labstack/echo/v4"

func (app *application) start() {
	e := echo.New()

	e.Static("/static", "ui/static")

	e.GET("/", app.home)
	e.GET("/user/signup", app.signup)
	e.GET("/user/login", app.login)

	e.Logger.Fatal(e.Start(app.addr))
}

package actions

import (
	"github.com/I1820/link/core"

	"github.com/labstack/echo/v4"
)

var linkApp *core.Application

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App(debug bool) *echo.Echo {
	e := echo.New()

	// prometheus middleware
	e.Use(NewPrometheusMiddleware("i1820_link"))

	// Routes
	e.GET("/about", AboutHandler)
	api := e.Group("/api")
	{
		mr := ModelsResource{}
		e.Resource("/models", mr)
		e.POST("/send", SendHandler)
	}

	return e
}

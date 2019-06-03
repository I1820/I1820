package actions

import (
	"github.com/labstack/echo/v4"
)

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
		mh := ModelsHandler{}
		api.GET("/models", mh.List)
		api.POST("/models", mh.Create)

		api.POST("/send", SendHandler)
		api.POST("/sendraw", sendRawHandler)
	}

	return e
}

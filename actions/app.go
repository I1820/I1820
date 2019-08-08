package actions

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App(debug bool) *echo.Echo {
	e := echo.New()

	// prometheus middleware
	e.Use(NewPrometheusMiddleware("i1820_link"))
	// prometheus metrics endpoint
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

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

package router

import (
	"github.com/I1820/I1820/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// App creates new instance of Echo and configures it
func App() *echo.Echo {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Pre(middleware.RemoveTrailingSlash())

	// prometheus middleware
	app.Use(NewPrometheusMiddleware(config.Namespace))

	// prometheus metrics endpoint
	app.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return app
}

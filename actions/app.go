package actions

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mongodb/mongo-go-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
)

// App creates new instance of Echo and configures it
func App(debug bool, databaseURL string) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Pre(middleware.RemoveTrailingSlash())

	if debug {
		app.Logger.SetLevel(log.DEBUG)
	}

	// validator
	app.Validator = &DefaultValidator{validator.New()}

	// prometheus middleware
	app.Use(NewPrometheusMiddleware("i1820_tm"))

	// routes
	app.GET("/about", AboutHandler)
	api := app.Group("/api")
	{
		tr := ThingsHandler{
			db: connectToDatabase(databaseURL),
		}

		pg := api.Group("/projects/:project_id")
		{
			pg.GET("/things", tr.List)
			pg.POST("/things", tr.Create)
		}
		api.DELETE("/things/:thing_id", tr.Destroy)
		api.GET("/things/:thing_id", tr.Show)
		api.PUT("/things/:thing_id", tr.Update)
		api.GET("/things/:thing_id/:t:(?:activate|deactivate)", tr.Activation)
	}

	return app
}

func connectToDatabase(url string) *mongo.Database {
	// Create mongodb connection
	client, err := mongo.NewClient(url)
	if err != nil {
		log.Fatalf("DB new client error: %s", err)
	}
	if err := client.Connect(context.Background()); err != nil {
		log.Fatalf("DB connection error: %s", err)
	}
	return client.Database("i1820")
}

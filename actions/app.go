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
func App(databaseURL string, debug bool) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Pre(middleware.RemoveTrailingSlash())

	if debug {
		app.Logger.SetLevel(log.DEBUG)
	}

	// Validator
	app.Validator = &DefaultValidator{validator.New()}

	// Routes
	app.GET("/about", AboutHandler)
	api := app.Group("/api")
	{
		pt := api.Group("/projects/:project_id/things/:thing_id")
		{
			qh := QueriesHandler{
				db: connectToDatabase(databaseURL),
			}
			pt.GET("/queries/list", qh.List)
			pt.POST("/queries/recently", qh.Recently)
			pt.POST("/queries/fetch", qh.Fetch)
			pt.POST("/queries/pfetch", qh.PartialFetch)
		}
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

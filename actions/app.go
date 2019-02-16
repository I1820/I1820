package actions

import (
	"context"

	"github.com/I1820/tm/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mongodb/mongo-go-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
)

// App creates new instance of Echo and configures it
func App() *echo.Echo {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Pre(middleware.RemoveTrailingSlash())

	if config.GetConfig().Debug {
		app.Logger.SetLevel(log.DEBUG)
	}

	// Validator
	app.Validator = &DefaultValidator{validator.New()}

	// Routes
	app.GET("/about", AboutHandler)
	api := app.Group("/api")
	{
		pg := api.Group("/projects/:project_id")
		{
			// /projects/{project_id}/things
			tr := ThingsHandler{
				db: connectToDatabase(),
			}
			pg.GET("/things", tr.List)
			pg.POST("/things", tr.Create)
			pg.DELETE("/things", tr.Destroy)
			pg.GET("/things/:thing_id", tr.Show)
			pg.PUT("/things/:thing_id", tr.Update)
			pg.POST("/things/geo", tr.GeoWithin)
			pg.POST("/things/tags", tr.HaveTags)
			pg.GET("/things/:thing_id/:t:(?:activate|deactivate)", tr.Activation)
		}
	}

	return app
}

func connectToDatabase() *mongo.Database {
	// Create mongodb connection
	url := config.GetConfig().Database.URL
	client, err := mongo.NewClient(url)
	if err != nil {
		log.Fatalf("DB new client error: %s", err)
	}
	if err := client.Connect(context.Background()); err != nil {
		log.Fatalf("DB connection error: %s", err)
	}
	return client.Database("i1820")
}

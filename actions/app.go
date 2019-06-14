package actions

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	// create mongodb connection
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		logrus.Fatalf("db new client error: %s", err)
	}
	// connect to the mongodb (change database here!)
	ctxc, donec := context.WithTimeout(context.Background(), 10*time.Second)
	defer donec()
	if err := client.Connect(ctxc); err != nil {
		logrus.Fatalf("db connection error: %s", err)
	}
	// is the mongo really there?
	ctxp, donep := context.WithTimeout(context.Background(), 2*time.Second)
	defer donep()
	if err := client.Ping(ctxp, readpref.Primary()); err != nil {
		logrus.Fatalf("db ping error: %s", err)
	}
	return client.Database("i1820")
}

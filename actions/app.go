package actions

import (
	"context"
	"reflect"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	// validator
	app.Validator = &DefaultValidator{validator.New()}

	// routes
	app.GET("/about", AboutHandler)
	api := app.Group("/api")
	{
		qh := QueriesHandler{
			db: connectToDatabase(databaseURL),
		}

		api.GET("/queries/projects/:project_id/list", qh.List)
		api.GET("/queries/things/:thing_id/parsed", qh.LastParsed)
		api.GET("/queries/things/:thing_id/fetcht", qh.FetchSingle)
		api.POST("/queries/fetch", qh.Fetch)
	}

	return app
}

func connectToDatabase(url string) *mongo.Database {
	rb := bson.NewRegistryBuilder()
	rb.RegisterTypeMapEntry(bsontype.EmbeddedDocument, reflect.TypeOf(bson.M{}))

	// create mongodb connection
	client, err := mongo.NewClient(options.Client().ApplyURI(url).SetRegistry(rb.Build()))
	if err != nil {
		log.Fatalf("db new client error: %s", err)
	}
	// connect to the mongodb (change database here!)
	ctxc, donec := context.WithTimeout(context.Background(), 10*time.Second)
	defer donec()
	if err := client.Connect(ctxc); err != nil {
		log.Fatalf("db connection error: %s", err)
	}
	// is the mongo really there?
	ctxp, donep := context.WithTimeout(context.Background(), 2*time.Second)
	defer donep()
	if err := client.Ping(ctxp, readpref.Primary()); err != nil {
		log.Fatalf("db ping error: %s", err)
	}
	return client.Database("i1820")
}

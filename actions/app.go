package actions

import (
	"context"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	mgo "github.com/mongodb/mongo-go-driver/mongo"

	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var db *mgo.Database

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_dm_session",
		})

		// If no content type is sent by the client
		// the application/json will be set, otherwise the client's
		// content type will be used.
		app.Use(contenttype.Add("application/json"))

		// Create mongodb connection
		url := envy.Get("DB_URL", "mongodb://172.18.0.1:27017")
		client, err := mgo.NewClient(url)
		if err != nil {
			buffalo.NewLogger("fatal").Fatalf("DB new client error: %s", err)
		}
		if err := client.Connect(context.Background()); err != nil {
			buffalo.NewLogger("fatal").Fatalf("DB connection error: %s", err)
		}
		db = client.Database("i1820")

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// Routes
		app.GET("/about", AboutHandler)
		api := app.Group("/api")
		{
			pt := api.Group("/projects/{project_id}/things/{thing_id}")
			{
				qr := QueriesResource{}
				pt.GET("/queries/list", qr.List)
				pt.POST("/queries/recently", qr.Recently)
				pt.POST("/queries/fetch", qr.Fetch)
				pt.POST("/queries/pfetch", qr.PartialFetch)
			}
		}
	}

	return app
}

package actions

import (
	"context"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var db *mgo.Database
var validate *validator.Validate

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
			SessionName: "_pm_session",
		})

		// If no content type is sent by the client
		// the application/json will be set, otherwise the client's
		// content type will be used.
		app.Use(contenttype.Add("application/json"))

		// create mongodb connection
		url := envy.Get("DB_URL", "mongodb://172.18.0.1:27017")
		client, err := mgo.NewClient(url)
		if err != nil {
			buffalo.NewLogger("fatal").Fatalf("DB new client error: %s", err)
		}
		if err := client.Connect(context.Background()); err != nil {
			buffalo.NewLogger("fatal").Fatalf("DB connection error: %s", err)
		}
		db = client.Database("i1820")

		// validator
		validate = validator.New()

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// prometheus collector
		rds := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "pm",
				Name:      "request_duration_seconds",
				Help:      "A histogram of latencies for requests.",
			},
			[]string{"path", "method", "code"},
		)

		prometheus.NewGoCollector()
		prometheus.MustRegister(rds)

		app.Use(func(next buffalo.Handler) buffalo.Handler {
			return func(c buffalo.Context) error {
				now := time.Now()

				defer func() {
					ws := c.Response().(*buffalo.Response)
					req := c.Request()

					rds.With(prometheus.Labels{
						"path":   req.URL.String(),
						"code":   strconv.Itoa(ws.Status),
						"method": req.Method,
					}).Observe(time.Since(now).Seconds())
				}()

				return next(c)
			}
		})

		// Routes
		app.GET("/about", AboutHandler)
		api := app.Group("/api")
		{

			// /projects
			pr := ProjectsResource{}
			api.Resource("/projects", pr)
			api.GET("/projects/{project_id}/logs", pr.Logs)

			// /projects/{project_id}/things
			pg := api.Group("/projects/{project_id}")
			{
				tr := ThingsResource{}
				pg.Resource("/things", tr)
				pg.POST("/things/geo", tr.GeoWithin)
				pg.GET("/things/{thing_id}/{t:(?:activate|deactivate)}", tr.Activation)
			}

			// /things/{thing_id}/assets
			api.Resource("/things/{thing_id}/assets", AssetsResource{})

			// /things/{thing_id}/connectivities
			api.Resource("/things/{thing_id}/connectivities", ConnectivitiesResource{})

			// /things/{thing_id}/tags
			tr := TagsResource{}
			api.POST("/things/{thing_id}/tags", tr.Create)
			api.GET("/things/{thing_id}/tags", tr.List)

			api.ANY("/runners/{project_id}/{path:.+}", RunnersHandler)
		}
		app.GET("/metrics", buffalo.WrapHandler(promhttp.Handler()))
	}

	return app
}

package actions

import (
	"context"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/unrolled/secure"

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
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))

		// Create mongodb connection
		url := envy.Get("DB_URL", "mongodb://172.18.0.1:27017")
		client, err := mgo.NewClient(url)
		if err != nil {
			buffalo.NewLogger("fatal").Fatalf("DB new client error: %s", err)
		}
		if err := client.Connect(context.Background()); err != nil {
			buffalo.NewLogger("fatal").Fatalf("DB connection error: %s", err)
		}
		db = client.Database("isrc")

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Collectors
		prometheus.Register(prometheus.NewGoCollector())

		rds := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "pm",
				Name:      "request_duration_seconds",
				Help:      "A histogram of latencies for requests.",
			},
			[]string{"path", "method", "code"},
		)
		prometheus.Register(rds)
		app.Use(func(h buffalo.Handler) buffalo.Handler {
			return func(c buffalo.Context) error {
				return buffalo.WrapHandler(promhttp.InstrumentHandlerDuration(
					rds.MustCurryWith(prometheus.Labels{"path": c.Value("current_path").(string)}),
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { h(c) }),
				))(c)
			}
		})

		// Routes
		app.GET("/about", AboutHandler)
		app.Mount("/", promhttp.Handler())
		api := app.Group("/api")
		{
			pr := ProjectsResource{}
			api.Resource("/projects", pr)
			api.GET("/projects/{project_id}/{t:(?:activate|deactivate)}", pr.Activation)
			api.GET("/projects/{project_id}/errors/{t:(?:lora|project)}", pr.Error)

			tr := ThingsResource{}
			api.Resource("/things", tr)
			api.GET("/things/{thing_id}/{t:(?:activate|deactivate)}", tr.Activation)
		}
	}

	return app
}

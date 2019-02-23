package actions

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"

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
				ch: connectToRabbitMQ(),
			}
			pg.GET("/things", tr.List)
			pg.POST("/things", tr.Create)
			pg.DELETE("/things/:thing_id", tr.Destroy)
			pg.GET("/things/:thing_id", tr.Show)
			pg.PUT("/things/:thing_id", tr.Update)
			pg.POST("/things/geo", tr.GeoWithin)
			pg.POST("/things/tags", tr.HaveTags)
			pg.GET("/things/:thing_id/:t:(?:activate|deactivate)", tr.Activation)

			// /projects/{project_id}/things/{thing_id}/assets
			ar := AssetsHandler{
				db: connectToDatabase(),
			}
			pg.GET("/things/:thing_id/assets", ar.List)
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

func connectToRabbitMQ() *amqp.Channel {
	// Makes a rabbitmq connection
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", config.GetConfig().Rabbit.User, config.GetConfig().Rabbit.Pass, config.GetConfig().Rabbit.Host))
	if err != nil {
		log.Fatalf("RabbitMQ connection error: %s", err)
	}

	// listen to rabbitmq close event
	go func() {
		for err := range conn.NotifyClose(make(chan *amqp.Error)) {
			log.Errorf("RabbitMQ connection is closed: %s", err)
			connectToRabbitMQ()
			return // connectToRabbitmQ will create new routine for error handling
		}
	}()

	// creates a rabbitmq channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ channel error: %s", err)
	}

	// direct exchange for thing creation and deletion events
	if err := ch.ExchangeDeclare(
		"i1820_things",
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	); err != nil {
		log.Fatalf("RabbitMQ failed to declare an exchange %s", err)
	}

	return ch
}

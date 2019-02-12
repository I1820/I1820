/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-08-2018
 * |
 * | File Name:     app.go
 * +===============================================
 */

package core

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	json "github.com/json-iterator/go"
	mgo "github.com/mongodb/mongo-go-driver/mongo"

	"github.com/I1820/dm/config"
	"github.com/I1820/types"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var instance *Application
var once sync.Once

// GetApplication returns global instance of core application.
func GetApplication() *Application {
	once.Do(func() {
		instance = new()
		instance.run()
	})
	return instance
}

// Application is a main part of dm component that consists of
// amqp client and protocols that provide information for amqp connectivity
// Application reads states from rabbitmq and stores them into database.
// Pipeline of application consists of following stages
// - Insert Stage
type Application struct {
	// rabbitmq connection
	stateConn *amqp.Connection
	stateChan *amqp.Channel

	Logger *logrus.Logger

	session *mgo.Client
	db      *mgo.Database

	// pipeline channels
	insertStream chan *types.State

	// in order to close the pipeline nicely
	// count number of stages so `Exit` can wait for all of them
	insertWG sync.WaitGroup

	IsRun bool
}

// New creates new application.
func new() *Application {
	a := Application{}

	a.Logger = logrus.New()

	// Create a mongodb connection
	url := config.GetConfig().Database.URL
	session, err := mgo.NewClient(url)
	if err != nil {
		a.Logger.Fatalf("Database client creation for %s error: %s", url, err)
	}
	a.session = session

	// pipeline channels
	a.insertStream = make(chan *types.State)

	return &a
}

// Run runs application. this function creates and connects amqp client.
func (a *Application) run() {
	// Makes a rabbitmq connection
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", config.GetConfig().Core.Broker.User, config.GetConfig().Core.Broker.Pass, config.GetConfig().Core.Broker.Host))
	if err != nil {
		a.Logger.Fatalf("RabbitMQ connection error: %s", err)
	}
	a.stateConn = conn

	// creates a rabbitmq channel
	ch, err := conn.Channel()
	if err != nil {
		a.Logger.Fatalf("RabbitMQ channel error: %s", err)
	}
	a.stateChan = ch

	// fanout exchange
	if err := a.stateChan.ExchangeDeclare(
		"i1820_fanout_states",
		"fanout",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	); err != nil {
		a.Logger.Fatalf("RabbitMQ failed to declare an exchange %s", err)
	}

	// listen to fanout exchange
	q, err := a.stateChan.QueueDeclare(
		"dm",  // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		a.Logger.Fatalf("RabbitMQ failed to declare a queue %s", err)
	}

	if err := a.stateChan.QueueBind(
		q.Name,                // queue name
		"",                    // routing key
		"i1820_fanout_states", // exchange
		false,
		nil,
	); err != nil {
		a.Logger.Fatalf("RabbitMQ failed to bind a queue %s", err)
	}

	msgs, err := a.stateChan.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		a.Logger.Fatalf("RabbitMQ failed to consume %s", err)
	}

	go func() {
		for msg := range msgs {
			a.consume(msg.Body)
		}
	}()

	// Connect to the mongodb
	if err := a.session.Connect(context.Background()); err != nil {
		a.Logger.Fatalf("DB connection error: %s", err)
	}
	a.db = a.session.Database("i1820")

	// pipeline stages
	for i := 0; i < runtime.NumCPU(); i++ {
		go a.insertStage()
		a.insertWG.Add(1)
	}

	a.IsRun = true
}

func (a *Application) consume(msg []byte) {
	var s types.State
	if err := json.Unmarshal(msg, &s); err != nil {
		a.Logger.WithFields(logrus.Fields{
			"component": "link",
		}).Errorf("Unmarshal data error: %s", err)
		return
	}
	a.insertStream <- &s
}

// Exit closes amqp connection then closes all channels and return from all pipeline stages
func (a *Application) Exit() {
	a.IsRun = false

	// close insert stream
	close(a.insertStream)
	a.insertWG.Wait()

	// disconnect from rabbitmq
	if err := a.stateChan.Close(); err != nil {
		a.Logger.Error(err)
	}
	if err := a.stateConn.Close(); err != nil {
		a.Logger.Error(err)
	}
}

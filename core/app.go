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
	log "github.com/sirupsen/logrus"

	"github.com/I1820/types"
	"github.com/streadway/amqp"
)

// Application is a main part of dm component that consists of
// amqp client and protocols that provide information for amqp connectivity
// Application reads states from rabbitmq and stores them into database.
// Pipeline of application consists of following stages
// - Insert Stage
type Application struct {
	// rabbitmq connection
	stateConn *amqp.Connection
	stateChan *amqp.Channel

	session *mgo.Client
	db      *mgo.Database

	// pipeline channels
	insertStream chan *types.State

	// in order to close the pipeline nicely
	// count number of stages so `Exit` can wait for all of them
	insertWG sync.WaitGroup

	// configuration parameters
	rabbitURL   string
	databaseURL string

	IsRun bool
}

// New creates new application.
func New(databaseURL string, rabbitURL string) *Application {
	return &Application{
		rabbitURL:   rabbitURL,
		databaseURL: databaseURL,

		insertStream: make(chan *types.State),
	}
}

// rabbitmqConnect connects to the rabbitmq and creates channel. it also provides
// a fail-safe way by reconnecting on connection failures.
func (a *Application) rabbitmqConnect() {
	// Makes a rabbitmq connection
	conn, err := amqp.Dial(a.rabbitURL)
	if err != nil {
		log.Fatalf("RabbitMQ connection error: %s", err)
	}
	a.stateConn = conn

	// listen to rabbitmq close event
	go func() {
		for err := range conn.NotifyClose(make(chan *amqp.Error)) {
			log.Errorf("RabbitMQ connection is closed: %s", err)
			a.rabbitmqConnect()
			return
		}
	}()

	// creates a rabbitmq channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ channel error: %s", err)
	}
	a.stateChan = ch
}

// Run runs application. this function creates and connects amqp client.
// this function also creates mongodb connection.
func (a *Application) Run() error {
	a.rabbitmqConnect()

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
		return fmt.Errorf("RabbitMQ failed to declare an exchange %s", err)
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
		return fmt.Errorf("RabbitMQ failed to declare a queue %s", err)
	}

	if err := a.stateChan.QueueBind(
		q.Name,                // queue name
		"",                    // routing key
		"i1820_fanout_states", // exchange
		false,
		nil,
	); err != nil {
		return fmt.Errorf("RabbitMQ failed to bind a queue %s", err)
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
		return fmt.Errorf("RabbitMQ failed to consume %s", err)
	}

	go func() {
		for msg := range msgs {
			a.consume(msg.Body)
		}
	}()

	// Create a mongodb connection
	session, err := mgo.NewClient(a.databaseURL)
	if err != nil {
		return fmt.Errorf("Database client creation for %s error: %s", a.databaseURL, err)
	}
	a.session = session

	// Connect to the mongodb
	if err := a.session.Connect(context.Background()); err != nil {
		return fmt.Errorf("DB connection error: %s", err)
	}
	a.db = a.session.Database("i1820")

	// pipeline stages
	for i := 0; i < runtime.NumCPU(); i++ {
		go a.insertStage()
		a.insertWG.Add(1)
	}

	a.IsRun = true
	return err
}

func (a *Application) consume(msg []byte) {
	var s types.State
	if err := json.Unmarshal(msg, &s); err != nil {
		log.WithFields(log.Fields{
			"component": "dm",
		}).Errorf("Unmarshal data error: %s", err)
		return
	}
	a.insertStream <- &s
}

// Exit closes amqp connection then closes all channels and return from all pipeline stages
func (a *Application) Exit() {
	a.IsRun = false

	// disconnect from rabbitmq
	if err := a.stateChan.Close(); err != nil {
		log.Error(err)
	}
	if err := a.stateConn.Close(); err != nil {
		log.Error(err)
	}

	// close insert stream
	close(a.insertStream)
	a.insertWG.Wait()
}

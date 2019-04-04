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
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
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
	conn *amqp.Connection
	ch   *amqp.Channel

	session *mongo.Client
	db      *mongo.Database

	// pipeline channels
	insertStream chan *types.State

	// in order to close the pipeline nicely
	// count number of stages so `Exit` can wait for all of them
	insertWG sync.WaitGroup

	IsRun bool

	// configuration parameters
	rabbitURL   string
	databaseURL string
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
func (a *Application) rabbitmqConnect() error {
	// Makes a rabbitmq connection
	conn, err := amqp.Dial(a.rabbitURL)
	if err != nil {
		return fmt.Errorf("rabbitmq connection error: %s", err)
	}
	a.conn = conn

	// listen to rabbitmq close event
	go func() {
		<-conn.NotifyClose(make(chan *amqp.Error))
		log.Errorf("RabbitMQ connection is closed: %s", err)
		if err := a.rabbitmqConnect(); err != nil {
			return // exits because there is no other way...
		}
		// exits because the new heath checker is coming.
	}()

	// creates a rabbitmq channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq channel error: %s", err)
	}
	a.ch = ch

	return nil
}

// Run runs application. this function creates and connects amqp client.
// this function also creates mongodb connection.
func (a *Application) Run() error {
	if err := a.rabbitmqConnect(); err != nil {
		return err
	}

	// fanout exchange of I1820 states
	// redefine of I1820 states exchange just for insurance
	if err := a.ch.ExchangeDeclare(
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

	// listen to fanout exchange of I1820 states
	qs, err := a.ch.QueueDeclare(
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

	if err := a.ch.QueueBind(
		qs.Name,               // queue name
		"",                    // routing key
		"i1820_fanout_states", // exchange
		false,
		nil,
	); err != nil {
		return fmt.Errorf("RabbitMQ failed to bind a queue %s", err)
	}

	go func() {
		msgs, err := a.ch.Consume(
			qs.Name, // queue
			"",      // consumer
			true,    // auto ack
			false,   // exclusive
			false,   // no local
			false,   // no wait
			nil,     // args
		)
		if err != nil {
			log.Errorf("RabbitMQ failed to consume %s", err)
		}

		for msg := range msgs {
			a.consume(msg.Body)
		}
	}()

	// direct exchange of I1820 things
	// redefine of I1820 things exchange just for insurance
	if err := a.ch.ExchangeDeclare(
		"i1820_things",
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("RabbitMQ failed to declare an exchange %s", err)
	}

	// listen to fanout exchange of I1820 things
	qtc, err := a.ch.QueueDeclare(
		"dm_thing_create", // name
		false,             // durable
		true,              // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return fmt.Errorf("RabbitMQ failed to declare a queue %s", err)
	}

	qtr, err := a.ch.QueueDeclare(
		"dm_thing_remove", // name
		false,             // durable
		true,              // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return fmt.Errorf("RabbitMQ failed to declare a queue %s", err)
	}

	if err := a.ch.QueueBind(
		qtc.Name,       // queue name
		"create",       // routing key
		"i1820_things", // exchange
		false,
		nil,
	); err != nil {
		return fmt.Errorf("RabbitMQ failed to bind a queue %s", err)
	}

	if err := a.ch.QueueBind(
		qtr.Name,       // queue name
		"remove",       // routing key
		"i1820_things", // exchange
		false,
		nil,
	); err != nil {
		return fmt.Errorf("RabbitMQ failed to bind a queue %s", err)
	}

	go func() {
		msgs, err := a.ch.Consume(
			qtc.Name, // queue
			"",       // consumer
			true,     // auto ack
			false,    // exclusive
			false,    // no local
			false,    // no wait
			nil,      // args
		)
		if err != nil {
			log.Errorf("RabbitMQ failed to consume %s", err)
		}

		for msg := range msgs {
			a.thingCreated(msg.Body)
		}
	}()

	go func() {
		msgs, err := a.ch.Consume(
			qtr.Name, // queue
			"",       // consumer
			true,     // auto ack
			false,    // exclusive
			false,    // no local
			false,    // no wait
			nil,      // args
		)
		if err != nil {
			log.Errorf("RabbitMQ failed to consume %s", err)
		}

		for msg := range msgs {
			a.thingRemoved(msg.Body)
		}
	}()

	// Create a mongodb connection
	session, err := mongo.NewClient(a.databaseURL)
	if err != nil {
		return fmt.Errorf("database client creation for %s error: %s", a.databaseURL, err)
	}
	a.session = session

	// Connect to the mongodb
	if err := a.session.Connect(context.Background()); err != nil {
		return fmt.Errorf("db connection error: %s", err)
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

func (a *Application) thingRemoved(msg []byte) {
	var t types.Thing
	if err := json.Unmarshal(msg, &t); err != nil {
		log.WithFields(log.Fields{
			"component": "dm",
		}).Errorf("Unmarshal data error: %s", err)
		return
	}

	// create data collection with following format
	// data.project_id.thing_id
	c := context.Background()
	cd := a.db.Collection(fmt.Sprintf("data.%s.%s", t.Project, t.ID))
	if err := cd.Drop(c); err != nil {
		// this error should not happen but in case of it happens you can ignore it safely.
		log.Errorf("Thing collection drop error: %s", err)
	}
	log.Infof("Thing %s collection is drop successfully", t.ID)
}

func (a *Application) thingCreated(msg []byte) {
	var t types.Thing
	if err := json.Unmarshal(msg, &t); err != nil {
		log.WithFields(log.Fields{
			"component": "dm",
		}).Errorf("Unmarshal data error: %s", err)
		return
	}

	// create data collection with following format
	// data.project_id.thing_id
	c := context.Background()
	cd := a.db.Collection(fmt.Sprintf("data.%s.%s", t.Project, t.ID))
	if _, err := cd.Indexes().CreateMany(
		c,
		[]mongo.IndexModel{
			{
				Keys: primitive.M{
					"at": 1,
				},
			},
			{
				Keys: primitive.M{
					"asset": 1,
				},
			},
		},
	); err != nil {
		// this error should not happen but in case of it happens you can ignore it safely.
		log.Errorf("Thing collection creation error: %s", err)
	}
	log.Infof("Thing %s collection is created successfully", t.ID)
}

// Exit closes amqp connection then closes all channels and return from all pipeline stages
func (a *Application) Exit() {
	a.IsRun = false

	// disconnect from rabbitmq
	if err := a.ch.Close(); err != nil {
		log.Error(err)
	}
	if err := a.conn.Close(); err != nil {
		log.Error(err)
	}

	// close insert stream
	close(a.insertStream)
	a.insertWG.Wait()
}

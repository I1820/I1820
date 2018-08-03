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

package app

import (
	"context"
	"log"
	"runtime"
	"time"

	pmclient "github.com/aiotrc/pm/client"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/gobuffalo/envy"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

// Application is a main component of uplink that consists of
// uplink protocols and mqtt client
type Application struct {
	cli paho.Client

	Logger *logrus.Logger

	protocols []Protocol

	pm pmclient.PM

	session *mgo.Client
	db      *mgo.Database

	// pipeline channels
	projectStream chan Data
	decodeStream  chan Data
	insertStream  chan Data
}

// Protocol is a uplink protocol
type Protocol interface {
	Topic() string
	Marshal([]byte) (Data, error)
}

// Data represents uplink data and metadata
type Data struct {
	Raw       []byte      // data before decode
	Data      interface{} // data after decode
	Timestamp time.Time   // when data received in uplink
	ThingID   string      // deveui
	RxInfo    interface{}
	TxInfo    interface{}
	Project   string // thing project identification
}

// New creates new application. this function creates mqtt client
func New() *Application {
	a := Application{}

	a.Logger = logrus.New()

	a.protocols = make([]Protocol, 0)

	// Create an MQTT client
	/*
		Port: 1883
		CleanSession: True
		Order: True
		KeepAlive: 30 (seconds)
		ConnectTimeout: 30 (seconds)
		MaxReconnectInterval 10 (minutes)
		AutoReconnect: True
	*/
	opts := paho.NewClientOptions()
	opts.AddBroker(envy.Get("BROKER_URL", "tcp://127.0.0.1:1883"))
	a.cli = paho.NewClient(opts)

	a.pm = pmclient.New(envy.Get("PM_URL", "http://127.0.0.1:8080"))

	// Create a mongodb connection
	url := envy.Get("DB_URL", "mongodb://127.0.0.1:27017")
	session, err := mgo.NewClient(url)
	if err != nil {
		a.Logger.Fatalf("DB new client error: %s", err)
	}
	a.session = session

	// pipeline channels
	a.projectStream = make(chan Data)
	a.decodeStream = make(chan Data)
	a.insertStream = make(chan Data)

	return &a
}

// Register registers protocol p on application a
func (a *Application) Register(p Protocol) {
	a.protocols = append(a.protocols, p)
}

// Run runs application. this function connects mqtt client and then register its topic
func (a *Application) Run() {
	// Connect to the MQTT Server.
	if t := a.cli.Connect(); t.Wait() && t.Error() != nil {
		a.Logger.Fatalf("MQTT session error: %s", t.Error())
	}

	// Connect to the mongodb
	if err := a.session.Connect(context.Background()); err != nil {
		log.Fatalf("DB connection error: %s", err)
	}
	// TODO db name must be configurable
	a.db = a.session.Database("isrc")

	// Subscribe to protocols topics
	for _, p := range a.protocols {
		if t := a.cli.Subscribe(p.Topic(), 0, a.mqttHandler(p)); t.Error() != nil {
			a.Logger.Fatalf("MQTT subscribe error: %s", t.Error())
		}
	}

	// pipeline stages
	for i := 0; i < runtime.NumCPU(); i++ {
		go a.project()
		go a.decode()
		go a.insert()
	}
}

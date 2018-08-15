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
	"fmt"
	"log"
	"runtime"

	pmclient "github.com/I1820/pm/client"
	"github.com/I1820/types"
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
	models    map[string]Model

	pm pmclient.PM

	session *mgo.Client
	db      *mgo.Database

	// pipeline channels
	projectStream chan types.Data
	decodeStream  chan types.Data
	insertStream  chan types.Data
}

// Protocol is a uplink/downlink protocol like lan or lora
type Protocol interface {
	TxTopic() string
	RxTopic() string

	Name() string

	Marshal([]byte) (types.Data, error)
}

// Model is a decoder/encoder interface like generic (based on user scripts) or aolab
type Model interface {
	Decode([]byte) interface{}
	Encode(interface{}) []byte

	Name() string
}

// New creates new application. this function creates mqtt client
func New() *Application {
	a := Application{}

	a.Logger = logrus.New()

	a.protocols = make([]Protocol, 0)
	a.models = make(map[string]Model)

	a.pm = pmclient.New(envy.Get("PM_URL", "http://127.0.0.1:8080"))

	// Create a mongodb connection
	url := envy.Get("DB_URL", "mongodb://127.0.0.1:27017")
	session, err := mgo.NewClient(url)
	if err != nil {
		a.Logger.Fatalf("DB new client error: %s", err)
	}
	a.session = session

	// pipeline channels
	a.projectStream = make(chan types.Data)
	a.decodeStream = make(chan types.Data)
	a.insertStream = make(chan types.Data)

	return &a
}

// RegisterProtocol registers protocol p on application a
func (a *Application) RegisterProtocol(p Protocol) {
	a.protocols = append(a.protocols, p)
}

// RegisterModel registers model m on application a
func (a *Application) RegisterModel(m Model) {
	a.models[m.Name()] = m
}

// Run runs application. this function connects mqtt client and then register its topic
func (a *Application) Run() {
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
	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(func(client paho.Client) {
		// Subscribe to protocols topics
		for _, p := range a.protocols {
			if t := a.cli.Subscribe(fmt.Sprintf("$share/i1820-link/%s", p.RxTopic()), 0, a.mqttHandler(p)); t.Error() != nil {
				a.Logger.Fatalf("MQTT subscribe error: %s", t.Error())
			}
		}
	})
	a.cli = paho.NewClient(opts)

	// Connect to the MQTT Server.
	if t := a.cli.Connect(); t.Wait() && t.Error() != nil {
		a.Logger.Fatalf("MQTT session error: %s", t.Error())
	}

	// Connect to the mongodb
	if err := a.session.Connect(context.Background()); err != nil {
		log.Fatalf("DB connection error: %s", err)
	}
	a.db = a.session.Database("i1820")

	// pipeline stages
	for i := 0; i < runtime.NumCPU(); i++ {
		go a.project()
		go a.decode()
		go a.insert()
	}
}

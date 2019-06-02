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
	"math/rand"
	"runtime"
	"time"

	"github.com/I1820/link/protocols"
	pmclient "github.com/I1820/pm/client"
	"github.com/I1820/types"
	paho "github.com/eclipse/paho.mqtt.golang"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Application is a main component of uplink that consists of
// uplink protocols and mqtt client
type Application struct {
	cli  paho.Client
	opts *paho.ClientOptions

	protocols []protocols.Protocol
	models    map[string]Model

	pm pmclient.PM

	session *mgo.Client
	db      *mgo.Database

	// pipeline channels
	projectStream chan types.Data
	decodeStream  chan types.Data
	insertStream  chan types.Data
}

// Model is a decoder/encoder interface like generic (based on user scripts) or aolab
type Model interface {
	Decode([]byte) interface{}
	Encode(interface{}) []byte

	Name() string
}

// New creates new application. this function creates mqtt client
func New(pmURL string, dbURL string, brokerURL string) *Application {
	a := Application{}

	a.protocols = make([]protocols.Protocol, 0)
	a.models = make(map[string]Model)

	// creates a pm communication link
	a.pm = pmclient.New(pmURL)

	// create a mongodb connection
	session, err := mgo.NewClient(dbURL)
	if err != nil {
		logrus.Fatalf("db new client error: %s", err)
	}
	a.session = session

	a.opts = paho.NewClientOptions()
	a.opts.AddBroker(brokerURL)
	a.opts.SetClientID(fmt.Sprintf("I1820-link-%d", rand.Intn(1024)))
	a.opts.SetOrderMatters(false)

	// pipeline channels
	a.projectStream = make(chan types.Data)
	a.decodeStream = make(chan types.Data)
	a.insertStream = make(chan types.Data)

	return &a
}

// RegisterProtocol registers protocol p on application a
func (a *Application) RegisterProtocol(p protocols.Protocol) {
	a.protocols = append(a.protocols, p)
}

// Protocols returns list of registered protocol's names
func (a *Application) Protocols() []string {
	names := make([]string, len(a.protocols))

	for i, p := range a.protocols {
		names[i] = p.Name()
	}

	return names
}

// RegisterModel registers model m on application a
func (a *Application) RegisterModel(m Model) {
	a.models[m.Name()] = m
}

// Models returns list of registered model's names
func (a *Application) Models() []string {
	names := make([]string, len(a.models))

	i := 0
	for n := range a.models {
		names[i] = n
		i++
	}

	return names
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
	a.opts.SetOnConnectHandler(func(client paho.Client) {
		// Subscribe to protocols topics
		for _, p := range a.protocols {
			if t := a.cli.Subscribe(fmt.Sprintf("%s", p.RxTopic()), 0, a.mqttHandler(p)); t.Error() != nil {
				logrus.Fatalf("mqtt subscribe error: %s", t.Error())
			}
		}
	})
	a.cli = paho.NewClient(a.opts)

	// Connect to the MQTT Server.
	if t := a.cli.Connect(); t.Wait() && t.Error() != nil {
		logrus.Fatalf("mqtt session error: %s", t.Error())
	}

	// Connect to the mongodb
	if err := a.session.Connect(context.Background()); err != nil {
		logrus.Fatalf("db connection error: %s", err)
	}
	a.db = a.session.Database("i1820")

	// pipeline stages
	for i := 0; i < runtime.NumCPU(); i++ {
		go a.project()
		go a.decode()
		go a.insert()
	}
}

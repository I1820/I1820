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
	"sync"
	"time"

	"github.com/I1820/link/models"
	"github.com/I1820/link/protocols"
	pmclient "github.com/I1820/pm/client"
	"github.com/I1820/types"
	paho "github.com/eclipse/paho.mqtt.golang"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

// Application is a main component of uplink that consists of
// uplink protocols and mqtt client
type Application struct {
	// mqtt configuration
	cli  paho.Client
	opts *paho.ClientOptions

	// models and protocols
	protocols []protocols.Protocol
	models    map[string]models.Model

	// pm connection
	pm pmclient.PM

	// database connection
	session *mgo.Client
	db      *mgo.Database

	// is core application running?
	IsRun bool

	// in order to close the pipeline nicely
	// count number of stages so `Exit` can wait for all of them
	projectWG sync.WaitGroup
	decodeWG  sync.WaitGroup
	insertWG  sync.WaitGroup

	// pipeline channels
	projectStream chan types.Data
	decodeStream  chan types.Data
	insertStream  chan types.Data
}

// New creates new application. this function creates mqtt client
func New(pmURL string, dbURL string, brokerURL string) (*Application, error) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	a := Application{}

	a.protocols = make([]protocols.Protocol, 0)
	a.models = make(map[string]models.Model)

	// creates a pm communication link
	a.pm = pmclient.New(pmURL)

	// create a mongodb connection
	session, err := mgo.NewClient(dbURL)
	if err != nil {
		return nil, err
	}
	a.session = session

	a.opts = paho.NewClientOptions()
	a.opts.AddBroker(brokerURL)
	a.opts.SetClientID(fmt.Sprintf("I1820-link-%d", rnd.Intn(1024)))
	a.opts.SetOrderMatters(false)

	// pipeline channels
	a.projectStream = make(chan types.Data)
	a.decodeStream = make(chan types.Data)
	a.insertStream = make(chan types.Data)

	return &a, nil
}

// RegisterProtocol registers protocol p on application a
func (a *Application) RegisterProtocol(p protocols.Protocol) {
	if a.IsRun {
		return // there is no way to add protocols in runtime
	}
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
func (a *Application) RegisterModel(m models.Model) {
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
func (a *Application) Run() error {
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
		// subscribe to protocols topics
		for _, p := range a.protocols {
			// use $share/i1820-link in front of the mqtt topic for group subscription
			if t := a.cli.Subscribe(fmt.Sprintf("%s", p.RxTopic()), 0, a.mqttHandler(p)); t.Error() != nil {
				logrus.Fatalf("mqtt subscribe error: %s", t.Error())
			}
		}
	})
	a.cli = paho.NewClient(a.opts)

	// connect to the MQTT Server
	if t := a.cli.Connect(); t.Wait() && t.Error() != nil {
		return t.Error()
	}

	// connect to the mongodb (change database here!)
	if err := a.session.Connect(context.Background()); err != nil {
		return err
	}
	a.db = a.session.Database("i1820")

	// pipeline stages
	for i := 0; i < runtime.NumCPU(); i++ {
		go a.project()
		a.projectWG.Add(1)
		go a.decode()
		a.decodeWG.Add(1)
		go a.insert()
		a.insertWG.Add(1)

	}

	a.IsRun = true

	return nil
}

// Exit closes amqp connection then closes all channels and return from all pipeline stages
func (a *Application) Exit() {
	a.IsRun = false

	// close project stream
	close(a.projectStream)
	a.projectWG.Wait()

	// close decode stream
	close(a.decodeStream)
	a.decodeWG.Wait()

	// close insert stream
	close(a.insertStream)
	a.insertWG.Wait()

	// disconnect aftere 10ms
	a.cli.Disconnect(10)
}

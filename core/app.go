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
	"runtime"
	"sync"

	types "github.com/I1820/I1820/model"
	"github.com/I1820/I1820/pkg/model"
	"github.com/I1820/I1820/rabbitmq"
	"github.com/I1820/I1820/store"
	"github.com/I1820/tm/client"
)

// Application is a main component of uplink that consists of
// uplink protocol and mqtt client
type Application struct {
	// model and protocol
	models map[string]model.Model

	// tm connection
	TMService client.TMService

	// data store
	Store *store.Data

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

	// RabbitMQ Producer
	rawProducer    *rabbitmq.Producer
	parsedProducer *rabbitmq.Producer
}

// New creates new application.
func New(tm client.TMService, st *store.Data, rpr, ppr *rabbitmq.Producer) *Application {
	return &Application{
		models:    make(map[string]model.Model),
		TMService: tm,
		Store:     st,

		projectStream: make(chan types.Data),
		decodeStream:  make(chan types.Data),
		insertStream:  make(chan types.Data),

		rawProducer:    rpr,
		parsedProducer: ppr,
	}
}

// RegisterModel registers model m on application a
func (a *Application) RegisterModel(m model.Model) {
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

// Handle handles given data into core pipeline
func (a *Application) Handle(d types.Data) {
	a.projectStream <- d
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
}

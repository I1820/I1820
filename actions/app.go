package actions

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"

	pmclient "github.com/aiotrc/pm/client"
	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// Config represents main configuration
var Config = struct {
	DB struct {
		URL string `default:"mongodb://127.0.0.1" env:"db_url"`
	}
	Broker struct {
		URL string `default:"127.0.0.1:1883" env:"broker_url"`
	}
	Decoder struct {
		Host string `default:"127.0.0.1" env:"decoder_host"`
	}
	PM struct {
		URL string `default:"http://127.0.0.1:8080" env:"pm_url"`
	}
	Redis struct {
		Addr string
	}
}{}

var db *mgo.Database
var pm pmclient.PM

// App creates configured mqtt application
func App() {
	// Load configuration
	if err := configor.Load(&Config, "config.yml"); err != nil {
		panic(err)
	}

	// Create a Mongo Session
	session, err := mgo.Connect(context.Background(), Config.DB.URL, nil)
	if err != nil {
		log.Fatalf("Mongo session %s: %v", Config.DB.URL, err)
	}

	// Create an MQTT client
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			log.WithFields(log.Fields{
				"component": "uplink",
			}).Errorf("MQTT Client %s", err)
		},
	})
	defer cli.Terminate()

	// Connect to the MQTT Server.
	if err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  Config.Broker.URL,
		ClientID: []byte(fmt.Sprintf("isrc-uplink-%d", rand.Int63())),
	}); err != nil {
		log.Fatalf("MQTT session %s: %s", Config.Broker.URL, err)
	}
	fmt.Printf("MQTT session %s has been created\n", Config.Broker.URL)

	// PM
	pm = pmclient.New(Config.PM.URL)

	// ISRC database
	db = session.Database("isrc")

	// Data collection
	indx, err := db.Collection("data").Indexes().CreateMany(
		context.Background(),
		[]mgo.IndexModel{
			mgo.IndexModel{
				Keys: bson.NewDocument(
					bson.EC.Int32("timestamp", 1),
				),
			},
			mgo.IndexModel{
				Keys: bson.NewDocument(
					bson.EC.Int32("thingid", 1),
					bson.EC.Int32("timestamp", 1),
				),
			},
			mgo.IndexModel{
				Keys: bson.NewDocument(
					bson.EC.String("data._location", "2dsphere"),
				),
			},
		},
	)
	if err != nil {
		log.Fatalf("Create index %v", err)
	}
	fmt.Printf("MongoDB \"data\" collection indexes: %v\n", indx)

	// Subscribe to topics
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			// https://docs.loraserver.io/use/getting-started/
			&client.SubReq{
				TopicFilter: []byte("application/+/device/+/error"),
				QoS:         mqtt.QoS0,
				Handler:     Error,
			},
			&client.SubReq{
				TopicFilter: []byte("application/+/device/+/rx"),
				QoS:         mqtt.QoS0,
				Handler:     Data,
			},
		},
	})
	if err != nil {
		log.Fatalf("MQTT subscription: %s", err)
	}
}

package actions

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/gobuffalo/envy"
	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"

	pmclient "github.com/aiotrc/pm/client"
	log "github.com/sirupsen/logrus"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

var db *mgo.Database
var pm pmclient.PM

// App creates configured mqtt application
func App() {
	// Create mongodb connection
	url := envy.Get("DB_URL", "mongodb://172.18.0.1:27017")
	session, err := mgo.NewClient(url)
	if err != nil {
		log.Fatalf("DB new client error: %s", err)
	}
	if err := session.Connect(context.Background()); err != nil {
		log.Fatalf("DB connection error: %s", err)
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
		Address:  envy.Get("BROKER_URL", "127.0.0.1:1883"),
		ClientID: []byte(fmt.Sprintf("isrc-uplink-%d", rand.Int63())),
	}); err != nil {
		log.Fatalf("MQTT session error: %s", err)
	}

	// PM
	pm = pmclient.New(envy.Get("PM_URL", "http://127.0.0.1:8080"))

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

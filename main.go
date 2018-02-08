/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 12-11-2017
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/aiotrc/uplink/decoder"
	"github.com/aiotrc/uplink/lora"
	"github.com/aiotrc/uplink/pm"
	"github.com/jinzhu/configor"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// Config represents main configuration
var Config = struct {
	DB struct {
		URL string `default:"127.0.0.1" env:"db_url"`
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
}{}

func main() {
	// Load configuration
	configor.Load(&Config, "config.yml")

	// Create a Mongo Session
	session, err := mgo.Dial(Config.DB.URL)
	if err != nil {
		log.Fatalf("Mongo session %s: %v", Config.DB.URL, err)
	}
	defer session.Close()
	fmt.Println("Mongo session has been created")

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// Raw collection
	cr := session.DB("isrc").C("raw")
	cr.Create(&mgo.CollectionInfo{
		Capped:  true,
		MaxDocs: 100,
	})

	// Create an MQTT client
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			log.Printf("MQTT client: %v", err)
		},
	})
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err = cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  Config.Broker.URL,
		ClientID: []byte(fmt.Sprintf("isrc-push-%d", rand.Int63())),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("MQTT session has been created")

	// PM
	pm := pm.New(Config.PM.URL)

	// Parsed collection
	cp := session.DB("isrc").C("parsed")
	cp.EnsureIndex(mgo.Index{
		Key: []string{"timestamp"},
	})

	// Subscribe to topics
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				// https://docs.loraserver.io/use/getting-started/
				TopicFilter: []byte("application/+/node/+/rx"),
				QoS:         mqtt.QoS0,
				Handler: func(topicName, message []byte) {
					var m lora.RxMessage
					if err := json.Unmarshal(message, &m); err != nil {
						log.Printf("Message: %v", err)
						return
					}
					fmt.Println(m)

					if err := cr.Insert(m); err != nil {
						log.Printf("Mongo insert [raw]: %v", err)
						return
					}

					// Find thing
					t, err := pm.GetThing(m.DevEUI)
					if err != nil {
						log.Printf("GetThing error: %s", err)
						return
					}
					// Create decoder
					decoder := decoder.New(fmt.Sprintf("http://%s:%s", Config.Decoder.Host, t.Project.Runner.Port))
					// Decode
					parsed, err := decoder.Decode(m.Data, m.DevEUI)
					if err != nil {
						log.Printf("Decoder: %v", err)
						return
					}

					var bdoc interface{}
					err = bson.UnmarshalJSON([]byte(parsed), &bdoc)
					if err != nil {
						log.Printf("Unmarshal JSON: %v\n %s", err, parsed)
						return
					}
					err = cp.Insert(&struct {
						Data      interface{}
						Timestamp time.Time
						ThingID   string
					}{
						Data:      bdoc,
						Timestamp: time.Now(),
						ThingID:   m.DevEUI,
					})
					if err != nil {
						log.Printf("Mongo insert [parsed]: %v", err)
						return
					}
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Wait for receiving a signal.
	<-sigc

	fmt.Println("18.20 As always ... left me alone")
}

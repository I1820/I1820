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
	"math/rand"
	"os"
	"os/signal"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	pmclient "github.com/aiotrc/pm/client"
	"github.com/aiotrc/uplink/decoder"
	"github.com/aiotrc/uplink/lora"
	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
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
	if err := configor.Load(&Config, "config.yml"); err != nil {
		panic(err)
	}

	// Create a Mongo Session
	session, err := mgo.Dial(Config.DB.URL)
	if err != nil {
		log.Fatalf("Mongo session %s: %v", Config.DB.URL, err)
	}
	defer session.Close()
	fmt.Printf("Mongo session %s has been created\n", Config.DB.URL)

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

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
	pm := pmclient.New(Config.PM.URL)

	// Parsed collection
	cp := session.DB("isrc").C("data")
	if err := cp.EnsureIndex(mgo.Index{
		Key: []string{"timestamp"},
	}); err != nil {
		panic(err)
	}

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
						log.WithFields(log.Fields{
							"component": "uplink",
						}).Errorf("JSON Unmarshal: %s", err)
						return
					}
					log.Info(m)

					// Find thing
					p, err := pm.GetThingProject(m.DevEUI)
					if err != nil {
						log.WithFields(log.Fields{
							"component": "uplink",
						}).Errorf("PM GetThingProject: %s", err)
						return
					}
					// TODO: thing activation
					/*
						if !t.Status {
							return
						}
					*/

					// Create decoder
					decoder := decoder.New(fmt.Sprintf("http://%s:%s", Config.Decoder.Host, p.Runner.Port))

					// Decode
					parsed, err := decoder.Decode(m.Data, m.DevEUI)
					if err != nil {
						log.WithFields(log.Fields{
							"component": "uplink",
						}).Errorf("Decode: %s", err)
						return
					}

					var bdoc interface{}
					if err := bson.UnmarshalJSON([]byte(parsed), &bdoc); err != nil {
						log.WithFields(log.Fields{
							"component": "uplink",
						}).Errorf("Unmarshal JSON: %s\n %q", err, parsed)
						return
					}

					defer func() {
						if err := cp.Insert(&struct {
							Raw       []byte
							Data      interface{}
							Timestamp time.Time
							ThingID   string
							RxInfo    []lora.RxInfo
							TxInfo    lora.TxInfo
							Project   string
						}{
							Raw:       m.Data,
							Data:      bdoc,
							Timestamp: time.Now(),
							ThingID:   m.DevEUI,
							RxInfo:    m.RxInfo,
							TxInfo:    m.TxInfo,
							Project:   p.Name,
						}); err != nil {
							log.WithFields(log.Fields{
								"component": "uplink",
							}).Errorf("Mongo insert: %s\n", err)
							return
						}
					}()
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("MQTT subscription: %s", err)
	}

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Wait for receiving a signal.
	<-sigc

	fmt.Println("18.20 As always ... left me alone")
}

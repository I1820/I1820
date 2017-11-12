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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

type raw struct {
	data []byte
}

func main() {
	// Create a Mongo Session
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// Raw collection
	cr := session.DB("isrc").C("parsed")

	// Create an MQTT client
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err = cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "127.0.0.1:1883",
		ClientID: []byte("isrc-push"),
	})
	if err != nil {
		panic(err)
	}

	// Subscribe to topics
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				// https://vernemq.com/docs/configuration/shared_subscriptions.html
				TopicFilter: []byte("$share/isrc/push"),
				QoS:         mqtt.QoS0,
				Handler: func(topicName, message []byte) {
					fmt.Println(string(topicName), string(message))
					cr.Insert(raw{
						[]byte("hello"),
					})
					fmt.Println("Decoding")
					r, _ := http.Post("http://127.0.0.1:8080/api/decode/me", "application/json", bytes.NewBuffer([]byte("hello")))
					b, _ := ioutil.ReadAll(r.Body)
					fmt.Println(string(b))
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// Parsed collection
	cp := session.DB("isrc").C("parsed")
	var bdoc interface{}
	err = bson.UnmarshalJSON([]byte(`{"id": 1,"name": "A green door","price": 12.50,"tags": ["home", "green"]}`), &bdoc)
	if err != nil {
		panic(err)
	}
	err = cp.Insert(&bdoc)

	if err != nil {
		panic(err)
	}

	for {
	}
}

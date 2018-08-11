/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 11-08-2018
 * |
 * | File Name:     main_downlink.go
 * +===============================================
 */

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/I1820/link/encoder"
	"github.com/I1820/link/lora"
	pmclient "github.com/I1820/pm/client"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// Config represents main configuration
var Config = struct {
	Broker struct {
		URL string `default:"127.0.0.1:1883" env:"broker_url"`
	}
	Encoder struct {
		Host string `default:"127.0.0.1" env:"encoder_host"`
	}
	PM struct {
		URL string `default:"http://127.0.0.1:8080" env:"pm_url"`
	}
	LanServer struct {
		URL string `default:"http://127.0.0.1:4000" env="lanserver_url"`
	}
}{}

var pm pmclient.PM
var cli *client.Client

func init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.WithFields(log.Fields{
		"component": "downlink",
	}).Writer()
}

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	r.Use(gin.ErrorLogger())

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)

		api.POST("/send", sendHandler)
		api.POST("/raw", sendRawHandler)
	}

	return r
}

func _main() {
	// Load configuration
	if err := configor.Load(&Config, "config.yml"); err != nil {
		panic(err)
	}

	pm = pmclient.New(Config.PM.URL)

	// Create an MQTT client
	cli = client.New(&client.Options{
		ErrorHandler: func(err error) {
			log.WithFields(log.Fields{
				"component": "downlink",
			}).Errorf("MQTT Client %s", err)
		},
	})
	defer cli.Terminate()

	// Connect to the MQTT Server.
	if err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  Config.Broker.URL,
		ClientID: []byte(fmt.Sprintf("isrc-push-%d", rand.Int63())),
	}); err != nil {
		log.Fatalf("MQTT session %s: %s", Config.Broker.URL, err)
	}
	fmt.Printf("MQTT session %s has been created\n", Config.Broker.URL)

	fmt.Println("Downlink AIoTRC @ 2018")

	r := handle()

	srv := &http.Server{
		Addr:    ":1373",
		Handler: r,
	}

	go func() {
		fmt.Printf("Downlink Listen: %s\n", srv.Addr)
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Listen Error:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Downlink Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown Error:", err)
	}
}

func aboutHandler(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")
}

func sendHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var r sendReq
	if err := c.BindJSON(&r); err != nil {
		return
	}

	p, err := pm.GetThingProject(r.ThingID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	encoder := encoder.New(fmt.Sprintf("http://%s:%s", Config.Encoder.Host, p.Runner.Port))

	raw, err := encoder.Encode(r.Data, r.ThingID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	b, err := json.Marshal(lora.TxMessage{
		Reference: "abcd1234",
		FPort:     r.FPort,
		Data:      raw,
		Confirmed: r.Confirmed,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := cli.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS0,
		TopicName: []byte(fmt.Sprintf("application/%s/node/%s/tx", r.ApplicationID, r.ThingID)),
		Message:   b,
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	lan, err := json.Marshal(struct {
		Data []byte
	}{
		Data: b,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	go func() {

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/devices/%s/push", Config.LanServer.URL, r.ThingID), bytes.NewReader(lan))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "aabbccddee11223344")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if resp.StatusCode != 200 {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Invalid lan server response"))
			return
		}
	}()

	c.JSON(http.StatusOK, raw)
}

func sendRawHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var r sendReq
	if err := c.BindJSON(&r); err != nil {
		return
	}

	b64, ok := r.Data.(string)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Invalid byte stream"))
		return
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	buffer := bytes.NewBuffer(raw)

	for raw := buffer.Next(r.SegmentSize); len(raw) != 0; raw = buffer.Next(r.SegmentSize) {
		log.Infof("Segment %v", raw)

		b, err := json.Marshal(lora.TxMessage{
			Reference: "abcd1234",
			FPort:     r.FPort,
			Data:      raw,
			Confirmed: r.Confirmed,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		for i := 0; i < r.RepeatNumber; i++ {
			if err := cli.Publish(&client.PublishOptions{
				QoS:       mqtt.QoS0,
				TopicName: []byte(fmt.Sprintf("application/%s/node/%s/tx", r.ApplicationID, r.ThingID)),
				Message:   b,
			}); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			log.Infof("MQTT Packet %s [%d]", b, i)

			time.Sleep(time.Duration(r.Sleep) * time.Second)
		}
	}

	c.JSON(http.StatusOK, raw)
}

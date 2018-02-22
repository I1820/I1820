/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-11-2017
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aiotrc/downlink/encoder"
	pmclient "github.com/aiotrc/pm/client"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
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
}{}

var pm pmclient.PM

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)

		api.POST("/send", sendHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	return r
}

func main() {
	// Load configuration
	if err := configor.Load(&Config, "config.yml"); err != nil {
		panic(err)
	}

	pm = pmclient.New(Config.PM.URL)

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
	var r sendReq

	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := pm.GetThing(r.ThingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	encoder := encoder.New(fmt.Sprintf("http://%s:%s", Config.Encoder.Host, t.Project.Runner.Port))

	raw, err := encoder.Encode(r.Data, r.ThingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", raw)

	fmt.Println(t)
}

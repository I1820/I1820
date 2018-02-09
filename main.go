/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 09-02-2018
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

	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
)

// Config represents main configuration
var Config = struct {
	DB struct {
		URL string `default:"127.0.0.1" env:"db_url"`
	}
}{}

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	return r
}

func main() {
	// Load configuration
	configor.Load(&Config, "config.yml")

	// Create a Mongo Session
	session, err := mgo.Dial(Config.DB.URL)
	if err != nil {
		log.Fatalf("Mongo session %s: %v", Config.DB.URL, err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	fmt.Println("DM AIoTRC @ 2017")

	srv := &http.Server{
		Addr:    ":1372",
		Handler: handle(),
	}

	go func() {
		fmt.Printf("DM Listen: %s\n", srv.Addr)
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
	fmt.Println("DM Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown Error:", err)
	}
}

func aboutHandler(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")
}

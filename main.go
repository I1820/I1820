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
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/aiotrc/dm/alias"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
)

// Config represents main configuration
var Config = struct {
	DB struct {
		URL string `default:"127.0.0.1" env:"db_url"`
	}
}{}

// ISRC database
var isrcDB *mgo.Database

// aliases
var aliases map[string]*alias.Alias

func init() {
	aliases = make(map[string]*alias.Alias)
}

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)

		api.GET("/things", thingsHandler)
		api.GET("/things/:thingid", thingDataHandler)
		api.POST("/things", thingsDataHandler)
		api.GET("/things/:thingid/key/:key", thingKeyDataHandler)

		api.POST("/alias", aliasCreateHandler)
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

	// Create a Mongo Session
	session, err := mgo.Dial(Config.DB.URL)
	if err != nil {
		log.Fatalf("Mongo session %s: %v", Config.DB.URL, err)
	}
	isrcDB = session.DB("isrc")
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

func thingsHandler(c *gin.Context) {
	var results []bson.M

	pipe := isrcDB.C("parsed").Pipe([]bson.M{
		{"$group": bson.M{"_id": "$thingid", "total": bson.M{"$sum": 1}}},
	})
	if err := pipe.All(&results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func thingKeyDataHandler(c *gin.Context) {
	var results []bson.M

	key := c.Param("key")
	id := c.Param("thingid")

	since, err := strconv.ParseInt(c.Query("since"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	until, err := strconv.ParseInt(c.Query("until"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := isrcDB.C("parsed").Find(bson.M{
		fmt.Sprintf("data.%s", key): bson.M{
			"$exists": true,
		},
		"thingid": id,
		"timestamp": bson.M{
			"$gt": time.Unix(since, 0),
			"$lt": time.Unix(until, 0),
		},
	}).Select(bson.M{
		"_id": false,
		fmt.Sprintf("data.%s", key): true,
		"timestamp":                 true,
		"thingid":                   true,
		"rxinfo":                    true,
	}).All(&results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, results)
}

func thingDataHandler(c *gin.Context) {
	var results []bson.M

	id := c.Param("thingid")

	since, err := strconv.ParseInt(c.Query("since"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	until, err := strconv.ParseInt(c.Query("until"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aliasing := bson.M{}
	if _, ok := aliases[id]; ok {
		for key, name := range aliases[id].Map {
			aliasing[name] = fmt.Sprintf("$data.%s", key)
		}
	}
	pipe := isrcDB.C("parsed").Pipe([]bson.M{
		{"$match": bson.M{
			"thingid": id,
			"timestamp": bson.M{
				"$gt": time.Unix(since, 0),
				"$lt": time.Unix(until, 0),
			},
		}},
		{"$addFields": bson.M{
			"info": func() interface{} {
				if len(aliasing) == 0 {
					return nil
				}
				return aliasing
			}(),
		}},
	})
	if err := pipe.All(&results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func thingsDataHandler(c *gin.Context) {
	var results []bson.M

	var json struct {
		ThingIDs []string `json:"thing_ids"`
		Since    int64
		Until    int64
	}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(json.ThingIDs) > 0 {
		pipe := isrcDB.C("parsed").Pipe([]bson.M{
			{"$match": bson.M{
				"thingid": bson.M{
					"$in": json.ThingIDs,
				},
				"timestamp": bson.M{
					"$gt": time.Unix(json.Since, 0),
					"$lt": time.Unix(json.Until, 0),
				},
			}},
			{"$sort": bson.M{"timestamp": 1}},
		})
		if err := pipe.All(&results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		pipe := isrcDB.C("parsed").Pipe([]bson.M{
			{"$match": bson.M{
				"thingid": bson.M{
					"$in": json.ThingIDs,
				},
				"timestamp": bson.M{
					"$gt": time.Unix(json.Since, 0),
					"$lt": time.Unix(json.Until, 0),
				},
			}},
			{"$sort": bson.M{"timestamp": 1}},
		})
		if err := pipe.All(&results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}

	c.JSON(http.StatusOK, results)

}

func aliasCreateHandler(c *gin.Context) {
	var json struct {
		Name string
		Map  map[string]string
	}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create, build and store
	a := alias.New(json.Name)
	for key, name := range json.Map {
		a.Add(key, name)
	}
	aliases[a.Name] = a

	c.JSON(http.StatusOK, a)
}

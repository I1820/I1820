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

	"github.com/aiotrc/dm/loraserver"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
)

// Config represents main configuration
var Config = struct {
	DB struct {
		URL string `default:"127.0.0.1" env:"db_url"`
	}

	LoRaServer struct {
		URL string `default:"platform.ceit.aut.ac.ir:50013" env:"loraserver_url"`
	}
}{}

// ISRC database
var isrcDB *mgo.Database

// ISRC loraserver.io
var isrcLoRaServer *loraserver.LoRaServer

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)

		api.GET("/things", thingsHandler)
		api.GET("/things/:thingid", thingDataHandler)
		api.GET("/things/:thingid/2", thingDataHandlerExperimental)
		api.POST("/things", thingsDataHandler)
		api.GET("/things/:thingid/key/:key", thingKeyDataHandler)

		api.GET("/gateway/:gatewayid/enable", gatewayLogEnable)
		api.GET("/gateway/:gatewayid", gatewayLogFetch)
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

	// Create a loraserver.io session
	l, err := loraserver.New(Config.LoRaServer.URL)
	if err != nil {
		log.Fatalf("loraserver.io session %s: %v", Config.LoRaServer.URL, err)
	}
	isrcLoRaServer = l

	// Create a Mongo Session
	session, err := mgo.Dial(Config.DB.URL)
	if err != nil {
		log.Fatalf("Mongo session %s: %v", Config.DB.URL, err)
	}
	isrcDB = session.DB("isrc")
	defer session.Close()

	// gateway collection
	cg := session.DB("isrc").C("gateway")
	if err := cg.EnsureIndex(mgo.Index{
		Key:         []string{"timestamp"},
		ExpireAfter: 10 * time.Second,
	}); err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	fmt.Println("DM AIoTRC @ 2018")

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

func thingDataHandlerExperimental(c *gin.Context) {
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

	pipe := isrcDB.C("parsed").Pipe([]bson.M{
		{"$match": bson.M{
			"thingid": id,
			"timestamp": bson.M{
				"$gt": time.Unix(since, 0),
				"$lt": time.Unix(until, 0),
			},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"dayOfYear": bson.M{"$dayOfYear": "$timestamp"},
				"hour":      bson.M{"$hour": "$timestamp"},
				"interval": bson.M{
					"$subtract": []bson.M{
						bson.M{"$minute": "$timestamp"},
						bson.M{"$mod": []interface{}{bson.M{"$minute": "$timestamp"}, 15}},
					},
				},
			},
			"count": bson.M{"$sum": 1},
			"last":  bson.M{"$last": "$data"},
		}},
		{"$sort": bson.M{"timestamp": -1}},
	})
	if err := pipe.All(&results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	pipe := isrcDB.C("parsed").Pipe([]bson.M{
		{"$match": bson.M{
			"thingid": id,
			"timestamp": bson.M{
				"$gt": time.Unix(since, 0),
				"$lt": time.Unix(until, 0),
			},
		}},
		{"$sort": bson.M{"timestamp": -1}},
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
			{"$sort": bson.M{"timestamp": -1}},
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
			{"$sort": bson.M{"timestamp": -1}},
		})
		if err := pipe.All(&results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}

	c.JSON(http.StatusOK, results)

}

func gatewayLogEnable(c *gin.Context) {
	mac := c.Param("gatewayid")

	ch, err := isrcLoRaServer.GatewayFrameStream(mac)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func() {
		for d := range ch {
			if err := isrcDB.C("gateway").Insert(d); err != nil {
				log.Println(err)
			}
		}
	}()

	c.JSON(http.StatusOK, true)
}

func gatewayLogFetch(c *gin.Context) {
	var results []bson.M = make([]bson.M, 0)

	mac := c.Param("gatewayid")

	skip, err := strconv.ParseInt(c.Query("skip"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pipe := isrcDB.C("gateway").Pipe([]bson.M{
		{"$match": bson.M{
			"mac": mac,
		}},
		{"$skip": skip},
		{"$limit": limit},
	})
	if err := pipe.All(&results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

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
	"strconv"
	"time"

	"github.com/aiotrc/pm/project"
	"github.com/aiotrc/pm/runner"
	"github.com/aiotrc/pm/thing"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/core/options"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
)

// Config represents main configuration
var Config = struct {
	DB struct {
		URL string `default:"mongodb://172.18.0.1:27017" env:"db_url"`
	}
}{}

// ISRC database
var isrcDB *mgo.Database

// init initiates global variables
func init() {
	// Load configuration
	if err := configor.Load(&Config, "config.yml"); err != nil {
		panic(err)
	}
}

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)

		api.POST("/project", projectNewHandler)
		api.GET("/project", projectListHandler)
		api.DELETE("/project/:name", projectRemoveHandler)
		api.POST("/project/:project/things", thingAddHandler)
		api.GET("/project/:project/logs", projectLogHandler)
		api.GET("/project/:project/", projectDetailHandler)

		api.GET("/things/:name", thingGetHandler)
		api.GET("/things/:name/activate", thingActivateHandler)
		api.GET("/things/:name/deactivate", thingDeactivateHandler)
		api.DELETE("/things/:name", thingRemoveHandler)
		api.GET("/things", thingListHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	return r
}

func setupDB() {
	// Create a Mongo Session
	client, err := mgo.Connect(context.Background(), Config.DB.URL, nil)
	if err != nil {
		log.Fatalf("Mongo session %s: %v", Config.DB.URL, err)
	}
	isrcDB = client.Database("isrc")
	// TODO removes logs on sepcific period (Time field)
}

func main() {
	setupDB()

	fmt.Println("PM AIoTRC @ 2018")

	r := handle()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		fmt.Printf("PM Listen: %s\n", srv.Addr)
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
	fmt.Println("PM Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown Error:", err)
	}
}

func aboutHandler(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")
}

func projectNewHandler(c *gin.Context) {
	var json projectReq
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := json.Name

	p, err := project.New(name, []runner.Env{
		{Name: "MONGO_URL", Value: Config.DB.URL},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := isrcDB.Collection("pm").InsertOne(context.Background(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func projectDetailHandler(c *gin.Context) {
	name := c.Param("project")

	var p project.Project

	dr := isrcDB.Collection("pm").FindOne(context.Background(), bson.NewDocument(
		bson.EC.String("name", name),
	))

	if err := dr.Decode(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func projectRemoveHandler(c *gin.Context) {
	name := c.Param("name")

	var p project.Project

	dr := isrcDB.Collection("pm").FindOne(context.Background(), bson.NewDocument(
		bson.EC.String("name", name),
	))

	if err := dr.Decode(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := p.Runner.Remove(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := isrcDB.Collection("pm").DeleteOne(context.Background(), bson.NewDocument(
		bson.EC.String("name", name),
	)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func projectListHandler(c *gin.Context) {
	ps := make([]project.Project, 0)

	cur, err := isrcDB.Collection("pm").Find(context.Background(), bson.NewDocument())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	for cur.Next(context.Background()) {
		var p project.Project

		if err := cur.Decode(&p); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ps = append(ps, p)
	}
	if err := cur.Close(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ps)
}

func thingAddHandler(c *gin.Context) {
	name := c.Param("project")

	var json thingReq
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t := thing.Thing{
		ID:     json.Name,
		Status: true,
	}

	dr := isrcDB.Collection("pm").FindOneAndUpdate(context.Background(), bson.NewDocument(
		bson.EC.String("name", name),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$push", bson.EC.Interface("things", t)),
	), mgo.Opt.ReturnDocument(options.After))

	var p project.Project

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Project %s not found", name)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func thingGetHandler(c *gin.Context) {
	name := c.Param("name")

	var p project.Project

	dr := isrcDB.Collection("pm").FindOne(context.Background(), bson.NewDocument(
		bson.EC.SubDocumentFromElements("things", bson.EC.SubDocumentFromElements("$elemMatch",
			bson.EC.String("id", name),
		)),
	))

	if err := dr.Decode(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func thingActivateHandler(c *gin.Context) {
	/*
		name := c.Param("name")

		if t, ok := things[name]; ok {
			t.Status = false
			c.JSON(http.StatusOK, t)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Thing %s not found", name)})
	*/
}

func thingDeactivateHandler(c *gin.Context) {
	/*
		name := c.Param("name")

		if t, ok := things[name]; ok {
			t.Status = true
			c.JSON(http.StatusOK, t)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Thing %s not found", name)})
	*/
}

func projectLogHandler(c *gin.Context) {
	var results = make([]*bson.Document, 0)

	id := c.Param("project")

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cur, err := isrcDB.Collection("errors").Aggregate(context.Background(), bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$match", bson.EC.String("project", id)),
		),
		bson.VC.DocumentFromElements(
			bson.EC.Int32("$limit", int32(limit)),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$sort", bson.EC.Int32("Time", -1)),
		),
	))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for cur.Next(context.Background()) {
		result := bson.NewDocument()

		if err := cur.Decode(result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		results = append(results, result)
	}
	if err := cur.Close(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func thingRemoveHandler(c *gin.Context) {
	name := c.Param("name")

	dr := isrcDB.Collection("pm").FindOneAndUpdate(context.Background(), bson.NewDocument(), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$pull", bson.EC.SubDocumentFromElements(
			"things", bson.EC.String("id", name)),
		),
	), mgo.Opt.ReturnDocument(options.After))

	var p project.Project

	if err := dr.Decode(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func thingListHandler(c *gin.Context) {
	results := make([]thing.Thing, 0)

	cur, err := isrcDB.Collection("pm").Aggregate(context.Background(), bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.String("$unwind", "$things"),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$replaceRoot", bson.EC.String("newRoot", "$things")),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$sort", bson.EC.Int32("Time", -1)),
		),
	))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for cur.Next(context.Background()) {
		var result thing.Thing

		if err := cur.Decode(&result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		results = append(results, result)
	}
	if err := cur.Close(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)

}

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
	"github.com/aiotrc/pm/thing"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Config represents main configuration
var Config = struct {
	DB struct {
		URL string `default:"172.18.0.1:27017" env:"db_url"`
	}
}{}

// in-memory databases for things and projects
var projects map[string]*project.Project
var things map[string]thing.Thing

// ISRC database
var isrcDB *mgo.Database

// init initiates global variables
func init() {
	projects = make(map[string]*project.Project)
	things = make(map[string]thing.Thing)
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

		api.GET("/things/:name", thingGetHandler)
		api.DELETE("/things/:name", thingRemoveHandler)
		api.GET("/things", thingListHandler)
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

	p, err := project.New(name, Config.DB.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	projects[name] = p
	c.JSON(http.StatusOK, p)
}

func projectRemoveHandler(c *gin.Context) {
	name := c.Param("name")

	if p, ok := projects[name]; ok {
		delete(projects, name)

		if err := p.Runner.Remove(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.String(http.StatusOK, name)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Project %s not found", name)})
}

func thingAddHandler(c *gin.Context) {
	project := c.Param("project")

	var json thingReq
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := json.Name

	if p, ok := projects[project]; ok {
		things[name] = thing.Thing{
			Project: p,
			ID:      name,
		}
		c.JSON(http.StatusOK, things[name])
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Project %s not found", name)})
}

func thingGetHandler(c *gin.Context) {
	name := c.Param("name")

	if t, ok := things[name]; ok {
		c.JSON(http.StatusOK, t)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Thing %s not found", name)})
}

func projectListHandler(c *gin.Context) {
	pl := make([]*project.Project, 0)

	for _, project := range projects {
		pl = append(pl, project)
	}

	c.JSON(http.StatusOK, pl)
}

func projectLogHandler(c *gin.Context) {
	var results []bson.M = make([]bson.M, 0)

	id := c.Param("project")

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := isrcDB.C("errors").Find(bson.M{
		"project": id,
	}).Limit(limit).All(&results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func thingRemoveHandler(c *gin.Context) {
	name := c.Param("name")

	if t, ok := things[name]; ok {
		c.JSON(http.StatusOK, t)
		delete(things, name)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Thing %s not found", name)})
}

func thingListHandler(c *gin.Context) {
	tl := make([]thing.Thing, 0)

	for _, thing := range things {
		tl = append(tl, thing)
	}

	c.JSON(http.StatusOK, tl)
}

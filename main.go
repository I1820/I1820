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
	"strings"
	"time"

	"github.com/aiotrc/pm/project"
	"github.com/aiotrc/pm/thing"
	"github.com/gin-gonic/gin"
)

// database or memory?
var projects map[string]*project.Project
var things map[string]thing.Thing

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
		api.DELETE("/project/:name", projectRemoveHandler)
		api.POST("/project/:project/things/", thingAddHandler)

		api.GET("/things/:name", thingGetHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	return r
}

func main() {
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

	name := strings.Replace(json.Name, " ", "_", -1)

	p, err := project.New(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	projects[name] = p
	c.JSON(http.StatusOK, p)
}

func projectRemoveHandler(c *gin.Context) {
	name := strings.Replace(c.Param("name"), " ", "_", -1)

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
	project := strings.Replace(c.Param("project"), " ", "_", -1)

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

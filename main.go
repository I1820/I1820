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

	"github.com/aiotrc/pm/project"
	"github.com/gin-gonic/gin"
)

var projects map[string]*project.Project

func init() {
	projects = make(map[string]*project.Project)
}

func main() {
	fmt.Println("PM by Parham Alvani")

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)
		api.GET("/project/:name", projectNewHandler)
		api.DELETE("/project/:name", projectRemoveHandler)
	}

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
	name := c.Param("name")
	p, err := project.New(name)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
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
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.String(http.StatusOK, name)
		return
	}
	c.String(http.StatusNotFound, name)
}

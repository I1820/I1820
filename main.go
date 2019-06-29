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
	"syscall"
	"time"

	"github.com/I1820/tm/db"
	"github.com/I1820/tm/handler"
	"github.com/I1820/tm/router"
	"github.com/I1820/tm/store"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("18.20 at Sep 07 2016 7:20 IR721")

	cfg := config()

	e := router.App(cfg.Debug, "i1820_tm")

	// routes
	db, err := db.New(cfg.Database.URL, "i1820")
	logrus.Fatal(err)

	th := handler.ThingsHandler{
		Store: store.Things{
			DB: db,
		},
	}

	e.GET("/about", handler.AboutHandler)
	api := e.Group("/api")
	{
		th.Register(api)
	}

	go func() {
		if err := e.Start(":1995"); err != http.ErrServerClosed {
			logrus.Fatalf("API Service failed with %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("18.20 As always ... left me alone")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("API Service failed on exit: %s", err)
	}
}

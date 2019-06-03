/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 12-11-2017
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

	"github.com/I1820/link/actions"
	"github.com/I1820/link/core"
	"github.com/I1820/link/models/aolab"
	"github.com/I1820/link/protocols/lan"
	"github.com/I1820/link/protocols/lora"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("18.20 at Sep 07 2016 7:20 IR721")

	// load configuration
	cfg := config()

	// creates the core application and registers the defaults
	core, err := core.New(cfg.PM.URL, cfg.Database.URL, cfg.Core.Broker.Addr)
	if err != nil {
		logrus.Fatal(err)
	}
	core.RegisterProtocol(lora.Protocol{})
	core.RegisterProtocol(lan.Protocol{})
	core.RegisterModel(aolab.Model{})
	if err := core.Run(); err != nil {
		logrus.Fatalf("Core Service failed with %s", err)
	}

	e := actions.App(cfg.Debug)
	go func() {
		if err := e.Start(":1373"); err != http.ErrServerClosed {
			log.Fatalf("API Service failed with %s", err)

		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("18.20 As always ... left me alone")

	core.Exit()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("API Service failed on exit: %s", err)
	}
}

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

	"github.com/I1820/dm/actions"
	"github.com/I1820/dm/core"
)

func main() {
	fmt.Println("18.20 at Sep 07 2016 7:20 IR721")

	e := actions.App()
	go func() {
		if err := e.Start(":1373"); err != http.ErrServerClosed {
			log.Fatalf("API Service failed with %s", err)

		}
	}()
	app := core.GetApplication()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("18.20 As always ... left me alone")

	app.Exit()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("API Service failed on exit: %s", err)
	}
}

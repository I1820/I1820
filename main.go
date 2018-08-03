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
	"fmt"
	"os"
	"os/signal"

	"github.com/aiotrc/uplink/app"
	"github.com/aiotrc/uplink/lan"
	"github.com/aiotrc/uplink/lora"
)

func main() {
	app := app.New()
	app.Register(lora.Protocol{})
	app.Register(lan.Protocol{})
	app.Run()
	fmt.Println("18.20 at Sep 07 2016 7:20 IR721")

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Wait for receiving a signal.
	<-sigc

	fmt.Println("18.20 As always ... left me alone")
}

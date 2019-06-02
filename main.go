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
	c := config()

	// creates the core application and registers the defaults
	linkApp, err := core.New(c.PM.URL, c.Database.URL, c.Core.Broker.Addr)
	if err != nil {
		logrus.Fatal(err)
	}
	linkApp.RegisterProtocol(lora.Protocol{})
	linkApp.RegisterProtocol(lan.Protocol{})
	linkApp.RegisterModel(aolab.Model{})
	if err := linkApp.Run(); err != nil {
		logrus.Fatal(err)
	}

	buffaloApp := actions.App()
	if err := buffaloApp.Serve(); err != nil {
		logrus.Fatal(err)
	}
}

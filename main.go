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
	"log"

	"github.com/I1820/link/actions"
	"github.com/I1820/link/aolab"
	"github.com/I1820/link/app"
	"github.com/I1820/link/lan"
	"github.com/I1820/link/lora"
)

func main() {
	linkApp := app.New()
	linkApp.RegisterProtocol(lora.Protocol{})
	linkApp.RegisterProtocol(lan.Protocol{})
	linkApp.RegisterModel(aolab.Model{})
	linkApp.Run()
	fmt.Println("18.20 at Sep 07 2016 7:20 IR721")

	buffaloApp := actions.App()
	if err := buffaloApp.Serve(); err != nil {
		log.Fatal(err)
	}
}

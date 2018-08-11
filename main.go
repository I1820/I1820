package main

import (
	"fmt"
	"log"

	"github.com/I1820/pm/actions"
)

func main() {
	fmt.Println("18.20 at Sep 07 2016 7:20 IR721")

	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}

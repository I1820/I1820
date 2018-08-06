package main

import (
	"log"

	"github.com/I1820/pm/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}

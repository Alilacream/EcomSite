package main

import (
	"log"
)

func main() {
	app := Setup()
	if err := app.run(app.routes()); err != nil {
		log.Fatal("Failed to start Server: ", err.Error())
	}
}

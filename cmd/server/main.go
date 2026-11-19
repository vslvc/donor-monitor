package main

import (
	"log"

	"donor-monitor/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/Espectro0/PlumbusBot/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(a); err != nil {
		log.Fatal(err)
	}
}

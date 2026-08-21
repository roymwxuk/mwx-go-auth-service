package main

import (
	"log"

	"github.com/roymwxuk/mwx-go-auth-service/config"
	"github.com/roymwxuk/mwx-go-auth-service/internal/app"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}

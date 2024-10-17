package main

import (
	"log"

	"github.com/aparnasukesh/payment-svc/config"
	"github.com/aparnasukesh/payment-svc/internals/di"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	server, err := di.InitResources(cfg)

	if err := server(); err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"log"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/config"
)

func main() {
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	server, diErr := di.InitializeAPI(config)
}

package main

import (
	"log"
	"github.com/rkachach/hss/cmd/config"
	"github.com/rkachach/hss/internal/api"
)

func main() {


	// Config initialization
	err := config.ReadConfig("config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	config.InitLogger()

	// Init API servers
	api.InitAPIRouter()
}

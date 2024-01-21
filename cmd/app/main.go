package main

import (
	"log"
	"github.com/rkachach/hss/cmd/config"
	"github.com/rkachach/hss/internal/api"
	"github.com/rkachach/hss/internal/dataStore"
)

// TODO find where to get information about dataStores
// we have also to see where the dataStore must be initialized
var store dataStore.DataStore = dataStore.OsFileSystem{}

func main() {

	// Config initialization
	err := config.ReadConfig("config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	config.InitLogger()

	store.Init(config.AppConfig.StoreConfig.Root)

	// Init API servers
	api.InitAPIRouter()
}

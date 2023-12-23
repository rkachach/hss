package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"fmt"
	//"github.com/rkachach/hss/internal/hss"
	"github.com/rkachach/hss/cmd/config"
	"github.com/rkachach/hss/internal/console"
)

const SlashSeparator string = "/"

func InitAPIRouter() {

	router := mux.NewRouter().SkipClean(true).UseEncodedPath()
	//apiRouter := router.PathPrefix(SlashSeparator).Subrouter()

	/////////////////////////////////////////////////
	////////////////// Management console operations
	/////////////////////////////////////////////////
	consoleRouter := http.NewServeMux()
	consoleRouter.HandleFunc("/config", console.ConsoleHandler)

	// listen on the console port
	go func() {
		addr := fmt.Sprintf(":%v", config.AppConfig.ConsolePort)
		log.Fatal(http.ListenAndServe(addr, consoleRouter))
	}()

	// listen on the main server port
	addr := fmt.Sprintf(":%v", config.AppConfig.ServerPort)
	log.Fatal(http.ListenAndServe(addr, router))
}

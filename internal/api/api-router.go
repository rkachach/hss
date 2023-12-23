package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"fmt"
	"github.com/rkachach/hss/internal/hss"
	"github.com/rkachach/hss/cmd/config"
	"github.com/rkachach/hss/internal/console"
)

const SlashSeparator string = "/"

func InitAPIRouter() {

	router := mux.NewRouter().SkipClean(true).UseEncodedPath()
	apiRouter := router.PathPrefix(SlashSeparator).Subrouter()

	var routers []*mux.Router
	routers = append(routers, apiRouter.PathPrefix("/{directory}").Subrouter())
	for _, router := range routers {

		////////////////////////////////////////
		////////////////// Directory operations
		////////////////////////////////////////

		// Bucket operations
		router.Methods(http.MethodPut).HandlerFunc(hss.Wrapper("PutDirectory", hss.PutDirectory))
		//router.Methods(http.MethodDelete).HandlerFunc(hss.Wrapper("DeleteDirectory", hss.DeleteDirectory))
		//router.Methods(http.MethodHead).HandlerFunc(hss.Wrapper("HeadDirectory", hss.HeadDirectory))
	}

	////////////////// Root operations
	//apiRouter.Methods(http.MethodGet).HandlerFunc(hss.Wrapper("ListDirectorys", hss.ListDirectorys))


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

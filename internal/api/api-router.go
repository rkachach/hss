package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"fmt"
	"github.com/rkachach/hss/internal/hss"
	"github.com/rkachach/hss/cmd/config"
	"github.com/rkachach/hss/internal/console"
	"github.com/gorilla/handlers"
)

const SlashSeparator string = "/"

func InitAPIRouter() {

	router := mux.NewRouter().SkipClean(true).UseEncodedPath()
	apiRouter := router.PathPrefix(SlashSeparator).Subrouter()

	var routers []*mux.Router
	routers = append(routers, apiRouter.PathPrefix("/{path:.+}").Subrouter())
	for _, router := range routers {

		////////////////////////////////////////
		////////////////// Directory operations
		////////////////////////////////////////




		// Directory operations
		router.Methods(http.MethodPost).HandlerFunc(hss.Wrapper("CreateDirectory", hss.CreateDirectory)).Queries("type", "directory")
		router.Methods(http.MethodGet).HandlerFunc(hss.Wrapper("ListDirectory", hss.ListDirectory)).Queries("type", "directory", "operation", "list")
		router.Methods(http.MethodGet).HandlerFunc(hss.Wrapper("GetDirectory", hss.GetDirectory)).Queries("type", "directory")
		router.Methods(http.MethodHead).HandlerFunc(hss.Wrapper("HeadDirectory", hss.HeadDirectory)).Queries("type", "directory")
		router.Methods(http.MethodDelete).HandlerFunc(hss.Wrapper("DeleteDirectory", hss.DeleteDirectory)).Queries("type", "directory")

		// File operations
		router.Methods(http.MethodPost).HandlerFunc(hss.Wrapper("CreateFile", hss.CreateFile)).Queries("type", "file")
		router.Methods(http.MethodGet).HandlerFunc(hss.Wrapper("GetFile", hss.GetFile)).Queries("type", "file")
		router.Methods(http.MethodHead).HandlerFunc(hss.Wrapper("HeadFile", hss.HeadFile)).Queries("type", "file")
		router.Methods(http.MethodDelete).HandlerFunc(hss.Wrapper("DeleteFile", hss.DeleteFile)).Queries("type", "file")
	}

	////////////////// Root operations
	apiRouter.Methods(http.MethodGet).HandlerFunc(hss.Wrapper("ListDirectory", hss.ListDirectory)).Queries("type", "directory", "operation", "list")

	// CORS middleware
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}), // Replace "*" with your allowed origins
	)

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
	log.Fatal(http.ListenAndServe(addr, cors(router)))
}

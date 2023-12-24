package hss

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rkachach/hss/internal/objectStore"
	"github.com/rkachach/hss/cmd/config"
)

func PutDirectory(w http.ResponseWriter, r *http.Request) {

	directoryName := mux.Vars(r)["directory"]

	err := objectStore.CreateDirectory(directoryName)
	if err != nil {
		// writeErrorResponse(w, errorCodes.ToAPIErr(ErrDirectoryAlreadyExists), r.URL)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteDirectory(w http.ResponseWriter, r *http.Request) {

	directoryName := mux.Vars(r)["directory"]
	objectStore.DeleteDirectory(directoryName)
	// Write success response.
	w.WriteHeader(http.StatusOK)
}

func HeadDirectory(w http.ResponseWriter, r *http.Request) {
	directoryName := mux.Vars(r)["directory"]
	dirInfo, err := objectStore.GetyDirectoryInfo(directoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Directory-Path", dirInfo.Path)
	w.Header().Set("Directory-Size", fmt.Sprintf("%v",dirInfo.Size))
	w.Header().Set("Directory-Files-Count", fmt.Sprintf("%v",dirInfo.FilesCount))
	w.WriteHeader(http.StatusOK)
}

func GetDirectory(w http.ResponseWriter, r *http.Request) {
	directoryName := mux.Vars(r)["directory"]

	dirInfo, err := objectStore.GetyDirectoryInfo(directoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("--> %v\n", dirInfo)

	// Convert directory info to a JSON response
	jsonResponse, err := json.Marshal(dirInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("err --> %v\n", err)
		return
	}

	fmt.Printf("json --> %v\n", jsonResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	count, err := w.Write(jsonResponse)
	if err != nil {
		fmt.Printf("err --> %v\n", err)
	} else {
		fmt.Printf("written --> %v\n", count)
	}
}

func ListDirectory(w http.ResponseWriter, r *http.Request) {
	directoryName := mux.Vars(r)["directory"]
	entries, err := objectStore.ListDirectory(directoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert file names to a JSON response
	jsonResponse, err := json.Marshal(entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func convertToMap(arr []string) map[string]bool {
	resultMap := make(map[string]bool)
	for _, key := range arr {
		resultMap[key] = true
	}
	return resultMap
}

func logRequestInfo(operation string, r *http.Request) {

	//TODO: see how can we generate and handle this map to avoid iterating over all the list
	selectedHeadersForLogging := convertToMap(config.AppConfig.Logging.SpecificHeaders)
	logMessage := fmt.Sprintf("========= Operation: %s\n", operation)
	logMessage += fmt.Sprintf("Path: %s\n", r.URL.Path)
	logMessage += fmt.Sprintf("ContentLengh: %v\n", r.ContentLength)

	// Log request query parameters
	logMessage += "Query parameters:\n"
	queryParams := r.URL.Query()
	for key, values := range queryParams {
		for _, value := range values {
			logMessage += fmt.Sprintf("- %s: %s\n", key, value)
		}
	}

	// Log request headers
	logMessage += "Received Headers:\n"
	for key, value := range r.Header {
		if len(selectedHeadersForLogging) > 0 {
			_, ok := selectedHeadersForLogging[key]
			if ok {
				logMessage += fmt.Sprintf("%s: %s\n", key, value)
			}
		} else {
			logMessage += fmt.Sprintf("- %s: %s\n", key, value)
		}
	}

	config.Logger.Println(logMessage)
}

func Wrapper(api string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequestInfo(api, r)
		f.ServeHTTP(w, r)
	}
}

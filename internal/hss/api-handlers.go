package hss

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rkachach/hss/internal/dataStore"
	"github.com/rkachach/hss/cmd/config"
	"io"
	"crypto/md5"
	"encoding/hex"
)

func PutDirectory(w http.ResponseWriter, r *http.Request) {

	dirPath := mux.Vars(r)["path"]
	err := dataStore.CreateDirectory(dirPath)
	if err != nil {
		// writeErrorResponse(w, errorCodes.ToAPIErr(ErrDirectoryAlreadyExists), r.URL)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteDirectory(w http.ResponseWriter, r *http.Request) {

	dirPath := mux.Vars(r)["path"]
	dataStore.DeleteDirectory(dirPath)
	// Write success response.
	w.WriteHeader(http.StatusOK)
}

func PutFile(w http.ResponseWriter, r *http.Request) {

	filePath := mux.Vars(r)["path"]
	dataStore.StartFileUpload(filePath)

	fileInfo, err := dataStore.StartFileUpload(filePath)
	if err == nil {
		filePartData, err := io.ReadAll(r.Body)
		if err != nil && err != io.EOF {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		fileInfo, err = dataStore.WriteFilePart(filePath, filePartData, 0)
		if err != nil {
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}

		FileBytes, err := dataStore.ReadFile(filePath)
		md5Hash := md5.Sum(FileBytes)
		md5Checksum := hex.EncodeToString(md5Hash[:])
		fileInfo.MD5sum = md5Checksum
		dataStore.UpdateFileInfo(filePath, fileInfo)


	} else {
		http.Error(w, "Error when creating a new upload", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", 0))
	w.WriteHeader(http.StatusOK)
}

func GetFile(w http.ResponseWriter, r *http.Request) {

	filePath := mux.Vars(r)["path"]
	fileBytes, err := dataStore.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Error reading file ", http.StatusNotFound)
		return
	}

	// Calculate MD5 checksum
	hash := md5.Sum(fileBytes)
	hashString := hex.EncodeToString(hash[:])
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-MD5", hashString)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileBytes)))

	//Copy the file content to the response writer
	_, err = w.Write(fileBytes)
	if err != nil {
		http.Error(w, "Error copying file content", http.StatusInternalServerError)
		return
	}
}

func HeadFile(w http.ResponseWriter, r *http.Request) {

}

func DeleteFile(w http.ResponseWriter, r *http.Request) {

}

func HeadDirectory(w http.ResponseWriter, r *http.Request) {
	dirPath := mux.Vars(r)["path"]
	dirInfo, err := dataStore.GetyDirectoryInfo(dirPath)
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
	dirPath := mux.Vars(r)["path"]
	dirInfo, err := dataStore.GetyDirectoryInfo(dirPath)
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
	dirPath := mux.Vars(r)["path"]
	entries, err := dataStore.ListDirectory(dirPath)
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

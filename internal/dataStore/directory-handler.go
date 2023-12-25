package dataStore

import (
	"sync"
	"path/filepath"
	"encoding/json"
	"os"
	"fmt"
	"time"
	"github.com/rkachach/hss/cmd/config"
	fsutils "github.com/rkachach/hss/internal/utils"
)

type DirectoryError struct {
	Op     string
	Key   string
	Err    error
}

func (e *DirectoryError) Error() string { return e.Op + " " + e.Key + ": " + e.Err.Error() }

type DirectoryInfo struct {
	// Name of the directory
	Name string `json:"name"`

	// Properties of the directory
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	FilesCount  int       `json:"files_count"`

	// Date and time when the directory was created/deleted
	CreatedTime time.Time `json:"created"`
	DeletedTime time.Time `json:"deleted,omitempty"`

	// Metadata for the directory
	Metadata map[string]string `json:"metadata,omitempty"`
}

var(
    lock sync.Mutex
)

func getDirectoryPath(dirPath string) string {
	dir := filepath.Dir(dirPath)
	return fmt.Sprintf("%s/%s", config.AppConfig.StoreConfig.Root, dir)
}

func getDirectoryInfoPath(dirPath string) string {
	dir := filepath.Dir(dirPath)
	return fmt.Sprintf("%s/%s/__info__.json", config.AppConfig.StoreConfig.Root, dir)
}

func writeDirectoryInfo(relativeDirPath string, directoryInfo *DirectoryInfo) {

	// Marshal struct to JSON
	jsonData, err := json.MarshalIndent(directoryInfo, "", "  ")
	if err != nil {
		fmt.Println("Error when writing directoryInfo:", err)
		return
	}

	// Write directory info file
	file, err := os.Create(getDirectoryInfoPath(relativeDirPath))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func CreateDirectory(relativeDirPath string, userMetadata map[string]string) error {

	lock.Lock()
	defer lock.Unlock()

	var directoryInfo DirectoryInfo
	directoryInfo.Name = relativeDirPath
	directoryInfo.CreatedTime = time.Now()
	directoryInfo.Metadata = userMetadata

	// Create a directory if it doesn't exist
	dirPath := getDirectoryPath(relativeDirPath)
	exists, err := fsutils.DirectoryExists(dirPath)
	if exists {
		config.Logger.Printf("CreateDirectory: %v already exists", config.AppConfig.StoreConfig.Root)
		return &DirectoryError{Op: "already exists", Key: relativeDirPath}
	} else {
		config.Logger.Printf("Directory '%v' created successfully", relativeDirPath)
	}

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return &DirectoryError{Op: "Error creating directory", Key: relativeDirPath}
	}

	writeDirectoryInfo(relativeDirPath, &directoryInfo)
	return nil
}

func GetDirectoryInfo(relativeDirPath string) (DirectoryInfo, error) {

	dirPath := getDirectoryInfoPath(relativeDirPath)
	file, err := os.Open(dirPath)
	if err != nil {
		fmt.Println("Error: directory path doesn't exist", dirPath)
		return DirectoryInfo{}, err
	}
	defer file.Close()

	// Decode JSON from the file
	var bucketInfo DirectoryInfo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&bucketInfo)
	if err != nil {
		fmt.Println("Error:", err)
		return DirectoryInfo{}, err
	}

	return bucketInfo, nil
}

func DeleteDirectory(relativeDirPath string) error {

	lock.Lock()
	defer lock.Unlock()

	dirPath := getDirectoryPath(relativeDirPath)
	exists, err := fsutils.DirectoryExists(dirPath)
	if !exists {
		config.Logger.Printf("getDirectory: Directory not found")
		return &DirectoryError{Op: "Error creating directory", Key: relativeDirPath}
	}

	// delete the directory from the data store
	err = os.RemoveAll(dirPath)
	if err != nil {
		fmt.Println("Error deleting directory:", err)
		//http.Error(w, "Error deleting directory", http.StatusNotFound)
		return &DirectoryError{Op: "Error deleting directory", Key: relativeDirPath}
	} else {
		fmt.Printf("Directory '%v' deleted successfully\n", relativeDirPath)
	}

	return nil
}

func ListDirectory(relativeDirPath string) ([]string, error) {
	dirPath := getDirectoryPath(relativeDirPath)
	fmt.Printf("Listing '%v' directory\n", dirPath)
	return fsutils.ListDirectoryEntries(dirPath)
}

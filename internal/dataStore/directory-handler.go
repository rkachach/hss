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

type ObjectError struct {
	Op     string
	Key   string
	Err    error
}

func (e *ObjectError) Error() string { return e.Op + " " + e.Key + ": " + e.Err.Error() }

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

func getDirectoryPath(directoryName string) string {
	dir := filepath.Dir(directoryName)
	return fmt.Sprintf("%s/%s", config.AppConfig.StoreConfig.Root, dir)
}

func getDirectoryInfoPath(directoryName string) string {
	dir := filepath.Dir(directoryName)
	return fmt.Sprintf("%s/%s/__info__.json", config.AppConfig.StoreConfig.Root, dir)
}

func writeDirectoryInfo(dirName string, directoryInfo *DirectoryInfo) {

	// Marshal struct to JSON
	jsonData, err := json.MarshalIndent(directoryInfo, "", "  ")
	if err != nil {
		fmt.Println("Error when writing directoryInfo:", err)
		return
	}

	// Write directory info file
	file, err := os.Create(getDirectoryInfoPath(dirName))
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

func CreateDirectory(directoryName string, userMetadata map[string]string) error {

	lock.Lock()
	defer lock.Unlock()

	var directoryInfo DirectoryInfo
	directoryInfo.Name = directoryName
	directoryInfo.CreatedTime = time.Now()
	directoryInfo.Metadata = userMetadata

	// Create a directory if it doesn't exist
	dirPath := getDirectoryPath(directoryName)
	exists, err := fsutils.DirectoryExists(dirPath)
	if exists {
		config.Logger.Printf("CreateDirectory: %v already exists", config.AppConfig.StoreConfig.Root)
		return &ObjectError{Op: "already exists", Key: directoryName}
	} else {
		config.Logger.Printf("Directory '%v' created successfully", directoryName)
	}

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return &ObjectError{Op: "Error creating directory", Key: directoryName}
	}

	writeDirectoryInfo(directoryName, &directoryInfo)
	return nil
}

func GetyDirectoryInfo(directoryName string) (DirectoryInfo, error) {

	// Create a directory if it doesn't exist
	dirPath := getDirectoryPath(directoryName)
	exists, err := fsutils.DirectoryExists(dirPath)
	if !exists {
		config.Logger.Printf("GetyDirectoryInfo: %v does not exist", dirPath)
		return DirectoryInfo{}, err
	}

	// Get directory info using a separate function
	var info DirectoryInfo
	info.Path = dirPath
	size, filesCount, creationTime, err := fsutils.GetDirectoryStats(dirPath)
	if err != nil {
		return info, err
	}

	info.Size = size
	info.FilesCount = filesCount
	info.CreatedTime = creationTime

	return info, nil
}

func DeleteDirectory(directoryName string) error {

	lock.Lock()
	defer lock.Unlock()

	dirPath := getDirectoryPath(directoryName)
	exists, err := fsutils.DirectoryExists(dirPath)
	if !exists {
		config.Logger.Printf("getDirectory: Directory not found")
		//http.NotFound(w, r)
		return &ObjectError{Op: "Error creating directory", Key: directoryName}
	}

	// delete the directory from the data store
	err = os.RemoveAll(dirPath)
	if err != nil {
		fmt.Println("Error deleting directory:", err)
		//http.Error(w, "Error deleting directory", http.StatusNotFound)
		return &ObjectError{Op: "Error deleting directory", Key: directoryName}
	} else {
		fmt.Printf("Directory '%v' deleted successfully\n", directoryName)
	}

	// Todo: delete all the directory objects recursively!
	return nil
}

func ListDirectory(directoryName string) ([]string, error) {
	dirPath := getDirectoryPath(directoryName)
	fmt.Printf("Listing '%v' directory\n", dirPath)
	return fsutils.ListDirectoryEntries(dirPath)
}

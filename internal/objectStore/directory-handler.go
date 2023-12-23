package objectStore

import (
	"sync"
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
	// Date and time when the directory was created/deleted
	Created time.Time `json:"created"`
	Deleted time.Time `json:"deleted,omitempty"`
	// Directory features
	Versioning    bool `json:"versioning"`
	ObjectLocking bool `json:"objectLocking"`
}

var(
    lock sync.Mutex
)

func getDirectoryPath(root, directoryName string) string {
	return fmt.Sprintf("%s/%s", root, directoryName)
}

func writeDirectoryInfo(directoryDir string, directoryInfo *DirectoryInfo) {

	// Marshal struct to JSON
	jsonData, err := json.MarshalIndent(directoryInfo, "", "  ")
	if err != nil {
		fmt.Println("Error when writing directoryInfo:", err)
		return
	}

	// Write directory info file
	file, err := os.Create(directoryDir + "/info.json")
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

func CreateDirectory(directoryName string) error {

	lock.Lock()
	defer lock.Unlock()

	var directoryInfo DirectoryInfo
	directoryInfo.Name = directoryName
	directoryInfo.Created = time.Now()

	// Create a directory if it doesn't exist
	directoryDir := getDirectoryPath(config.AppConfig.StoreConfig.Root, directoryName)
	exists, err := fsutils.DirectoryExists(directoryDir)
	if exists {
		config.Logger.Printf("CreateDirectory: %v already exists", config.AppConfig.StoreConfig.Root)
		return &ObjectError{Op: "already exists", Key: directoryName}
	} else {
		config.Logger.Printf("Directory '%v' created successfully", directoryName)
	}

	err = os.MkdirAll(directoryDir, 0755)
	if err != nil {
		return &ObjectError{Op: "Error creating directory", Key: directoryName}
	}

	writeDirectoryInfo(directoryDir, &directoryInfo)
	return nil
}

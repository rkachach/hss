package dataStore

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"github.com/google/uuid"
	"github.com/rkachach/hss/cmd/config"
	fsutils "github.com/rkachach/hss/internal/utils"
)

type DataStore interface {
  Init(dataStore string) error
  IsMetadataFile(filename string) bool
  StartFileUpload(filePath string, userMetadata map[string]string) (FileInfo, error)
  ReadFileInfo(filePath string) (FileInfo, error)
  WriteFilePart(filePath string, objectPartData []byte, PartNumber int) (FileInfo, error)
  ReadFile(filePath string) ([]byte, error) // TODO: this should be a FilePartReader or something like that
  DeleteFile(filePath string) error
  UpdateFileInfo(filePath string, fileInfo FileInfo) error
  CreateDirectory(relativeDirPath string, userMetadata map[string]string) error
  GetDirectoryInfo(relativeDirPath string) (DirectoryInfo, error)
  DeleteDirectory(relativeDirPath string) error
  ListDirectory(relativeDirPath string) ([]ElementExtendedInfo, error)
}

type OsFileSystem struct {
}

type FileError struct {
	Op     string
	Key   string
	Err    error
}

func (e *FileError) Error() string { return e.Op + " " + e.Key + ": " + e.Err.Error() }

type ElementExtendedInfo struct {
	Name         string    `json:"name"`
	Type          string   `json:"type"`
	Key          string    `json:"key"`
	LastModified time.Time `json:"lastModified"`
	Size         int64     `json:"size"`
}

type FileInfo struct {
	Name         string    `json:"name"`
	Key          string    `json:"key"`
	LastModified time.Time `json:"lastModified"`
	Size         int64     `json:"size"`
	Checksum     []byte    `json:"checksum"`
	UploadID     string    `json:"uploadID"`
	MD5sum       string    `json:"MD5sum"`

	// Metadata for the directory
	Metadata map[string]string `json:"metadata,omitempty"`
}

func (store OsFileSystem) Init(dataStore string) error {
	err := os.MkdirAll(dataStore, 0755)
	if err != nil {
		return &DirectoryError{Op: "Error creating directory", Key: dataStore}
	}
	return err
}

func (store OsFileSystem) IsMetadataFile(filename string) bool {
	return strings.HasSuffix(filename, ".json")
}

func getFilePath(filePath string) string {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)
	return fmt.Sprintf("%s/%s/%s", config.AppConfig.StoreConfig.Root, dir, filename)
}

func getFileInfoPath(filePath string) string {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)
	return fmt.Sprintf("%s/%s/__%s__.json", config.AppConfig.StoreConfig.Root, dir, filename)
}

func writeFileInfo(filePath string, fileInfo *FileInfo) error {

	jsonData, err := json.MarshalIndent(fileInfo, "", "  ")
	if err != nil {
		fmt.Println("Error when writing Fileinfo:", err)
		return err
	}

	infoFilePath := getFileInfoPath(filePath)
	file, err := os.Create(infoFilePath)
	if err != nil {
		fmt.Printf("os.Create Error: %v (path: %v)\n", err, infoFilePath)
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("file.Write Error:", err)
		return err
	}

	return nil
}

func (store OsFileSystem) ReadFileInfo(filePath string) (FileInfo, error) {

	file, err := os.Open(getFileInfoPath(filePath))
	if err != nil {
		fmt.Println("Error:", err)
		return FileInfo{}, err
	}
	defer file.Close()

	var fileInfo FileInfo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&fileInfo)
	if err != nil {
		fmt.Println("Error:", err)
		return FileInfo{}, err
	}

	return fileInfo, nil
}

// FIXME: fix createion on files like "/", this could go pretty wrogn if we create a file in host root instead of 
// "datastore root"
func (store OsFileSystem) StartFileUpload(filePath string, userMetadata map[string]string) (FileInfo, error){
	lock.Lock()
	defer lock.Unlock()

	_, err := store.ReadFileInfo(filePath)
	if err == nil {
		// File already exists
		fmt.Printf("File %v already exsits\n", filePath)
		return FileInfo{}, err
	}

	fileInfo := FileInfo{Name: filePath,
		Key: filePath,
		LastModified: time.Now().UTC(),
		UploadID: uuid.New().String(),
		Size: 0,
		Metadata: userMetadata}

	err = writeFileInfo(filePath, &fileInfo)
	return fileInfo, err
}

// TODO: will file parts arrive syncrhonously or async ? Treat unordered part arrivals
// TODO: write parts with buffering, there is too many opens. We will probably need:
//             - Open file
//             - Close file
//             - Read range
//
func (store OsFileSystem) WriteFilePart(filePath string, objectPartData []byte, PartNumber int) (FileInfo, error) {

	lock.Lock()
	defer lock.Unlock()

	fileInfo, err := store.ReadFileInfo(filePath)
	if err != nil {
		fmt.Printf("Error reading file info: %v\n", err)
		return FileInfo{}, err
	}

	fsutils.AppendToFile(getFilePath(filePath), objectPartData)

	// Update object info
	fileInfo.Size = fileInfo.Size + int64(len(objectPartData))
	fileInfo.LastModified = time.Now().UTC()
	err = writeFileInfo(filePath, &fileInfo)
	if err != nil {
		fmt.Printf("Error writing file info: %v\n", err)
		return FileInfo{}, err
	}

	return fileInfo, nil
}

func (store OsFileSystem) ReadFile(filePath string) ([]byte, error) {

	lock.Lock()
	defer lock.Unlock()

	file, err := os.Open(getFilePath(filePath))
	if err != nil {
		return nil, &FileError{Op: "Error reading object", Key: filePath}
	}
	defer file.Close()

	// Get file info to set content length
	_, err = file.Stat()
	if err != nil {
		return nil, &FileError{Op: "Error reading object", Key: filePath}
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (store OsFileSystem) UpdateFileInfo(filePath string, fileInfo FileInfo) error {

	lock.Lock()
	defer lock.Unlock()

	err := writeFileInfo(filePath, &fileInfo)
	if err != nil {
		return err
	}

	return nil
}

func (store OsFileSystem) DeleteFile(filePath string) error {

	lock.Lock()
	defer lock.Unlock()

	_, err := store.ReadFileInfo(filePath)
	if err != nil {
		config.Logger.Printf("Error When reading info file for %v: %v", filePath, err)
	}

	err = os.Remove(getFileInfoPath(filePath))
	if err != nil {
		config.Logger.Printf("Error When removing file for %v: %v", filePath, err)
	}

	fmt.Printf("Deleting file: %v\n", getFilePath(filePath))
	err = os.Remove(getFilePath(filePath))
	if err != nil {
		config.Logger.Printf("Error When removing %v", filePath)
		return &FileError{Op: "Error deleting object", Key: filePath}
	}

	return nil
}

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

func isSubdirectory(parent, child string) bool {
	relPath, err := filepath.Rel(parent, child)
	if err != nil {
		// Error occurred, indicating that the paths are not related
		return false
	}
	// Check if the relative path starts with ".." or "../"
	return relPath == "." || relPath == ".." || relPath[:3] == "../"
}

func getDirectoryPath(dirPath string) (string, error) {
	// Concatenate the directory path with the root directory path
	finalPath := filepath.Join(config.AppConfig.StoreConfig.Root, dirPath)

	// Check if the directory path contains upward directory traversal
	if !isSubdirectory(finalPath, config.AppConfig.StoreConfig.Root) {
		return "", fmt.Errorf("invalid directory path")
	}

	return finalPath, nil
}

func getDirectoryInfoPath(dirPath string) string {
	return fmt.Sprintf("%s/%s/__info__.json", config.AppConfig.StoreConfig.Root, dirPath)
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

func (store OsFileSystem) CreateDirectory(relativeDirPath string, userMetadata map[string]string) error {

	lock.Lock()
	defer lock.Unlock()

	var directoryInfo DirectoryInfo
	directoryInfo.Name = relativeDirPath
	directoryInfo.CreatedTime = time.Now()
	directoryInfo.Metadata = userMetadata

	// Create a directory if it doesn't exist
	dirPath, err := getDirectoryPath(relativeDirPath)
	if err != nil {
		config.Logger.Printf("CreateDirectory: %v invalid path", dirPath)
		return &DirectoryError{Op: "CreateDirectory", Err: err, Key: relativeDirPath}
	}

	exists, err := fsutils.DirectoryExists(dirPath)
	if exists {
		config.Logger.Printf("CreateDirectory: %v already exists", dirPath)
		return &DirectoryError{Op: "already exists", Key: relativeDirPath}
	}

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return &DirectoryError{Op: "Error creating directory", Key: relativeDirPath}
	}
	config.Logger.Printf("Directory '%v' created successfully", dirPath)

	writeDirectoryInfo(relativeDirPath, &directoryInfo)
	return nil
}

func fillDirectoryInfo(info map[string]interface{}) DirectoryInfo {
    dirInfo := DirectoryInfo{}
    dirInfo.Name = info["Name"].(string)
    dirInfo.Path = info["Path"].(string)
    dirInfo.Size = info["Size"].(int64)
    dirInfo.FilesCount = info["FilesCount"].(int)
    dirInfo.CreatedTime = info["CreatedTime"].(time.Time)
    return dirInfo
}

func (store OsFileSystem) GetDirectoryInfo(relativeDirPath string) (DirectoryInfo, error) {

	dirInfoPath := getDirectoryInfoPath(relativeDirPath)
	file, err := os.Open(dirInfoPath)
	if err == nil {
		var dirInfo DirectoryInfo
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&dirInfo)
		if err != nil {
			fmt.Println("Error:", err)
			return DirectoryInfo{}, err
		}
		return dirInfo, nil
	}

	// Directory info file doesn't exist, let's get dir info from filesytem
	dirPath, err := getDirectoryPath(relativeDirPath)
	if err != nil {
		config.Logger.Printf("Cannot get directory info: %v ", dirPath)
		return DirectoryInfo{}, &DirectoryError{Op: "GetDirecotryInfo", Err: err, Key: relativeDirPath}
	}

	dirInfoMap, err := fsutils.GetDirectoryInfo(dirPath)
	if err == nil {
		return fillDirectoryInfo(dirInfoMap), nil
	} else {
		return DirectoryInfo{}, &DirectoryError{Op: "GetDirecotryInfo", Err: err, Key: relativeDirPath}
	}
}

func (store OsFileSystem) DeleteDirectory(relativeDirPath string) error {

	lock.Lock()
	defer lock.Unlock()

	dirPath, err := getDirectoryPath(relativeDirPath)
	if err != nil {
		config.Logger.Printf("DeleteDirectory: %v invalid path", dirPath)
		return &DirectoryError{Op: "DeleteDirectory", Err: err, Key: relativeDirPath}
	}

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

func (store OsFileSystem) ListDirectory(relativeDirPath string) ([]ElementExtendedInfo, error) {
	dirPath, err := getDirectoryPath(relativeDirPath)
	if err != nil {
		config.Logger.Printf("Cannot list directory: %v ", dirPath)
		return nil, &DirectoryError{Op: "ListDirectory", Err: err, Key: relativeDirPath}
	}

	fmt.Printf("Listing '%v' directory\n", dirPath)
	var dirEntries []ElementExtendedInfo
	elements, err := fsutils.ListDirectoryWithDetails(dirPath)
	if err == nil {
		for _, entry := range elements {
			if ! store.IsMetadataFile(entry.Name){
				elementType := ""
				if entry.IsDirectory {
					elementType = "directory"
				} else {
					elementType = "file"
				}
				dirElement := ElementExtendedInfo{Name:entry.Name, Type: elementType}
				dirEntries = append(dirEntries, dirElement)
			}
		}
	}

	return dirEntries, err
}

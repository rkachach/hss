package dataStore

import(
	"path/filepath"
	"fmt"
	"io"
	"encoding/json"
	"os"
	"time"
	"strings"
	"github.com/rkachach/hss/cmd/config"
	fsutils "github.com/rkachach/hss/internal/utils"
	"github.com/google/uuid"
)

type FileError struct {
	Op     string
	Key   string
	Err    error
}

func (e *FileError) Error() string { return e.Op + " " + e.Key + ": " + e.Err.Error() }

type FileInfo struct {
	Name         string    `json:"name"`
	Key          string    `json:"key"`
	LastModified time.Time `json:"lastModified"`
	Size         int64     `json:"size"`
	Checksum     []byte    `json:"checksum"`
	UploadID     string    `json:"uploadID"`
	MD5sum       string    `json:"MD5sum"`
}

func IsMetadataFile(filename string) bool {
	return strings.HasSuffix(filename, ".json")
}

func getFilePath(filePath string) string {
	return fmt.Sprintf("%s/%s", config.AppConfig.StoreConfig.Root, filePath)
}

func getFileInfoPath(filePath string) string {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)
	return fmt.Sprintf("%s/%s/__%s.json", config.AppConfig.StoreConfig.Root, dir, filename)
}

func writeFiletInfo(filePath string, fileInfo *FileInfo) error {

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

func readFileInfo(filePath string) (FileInfo, error) {

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

func StartFileUpload(filePath string) (FileInfo, error){
	lock.Lock()
	defer lock.Unlock()

	_, err := readFileInfo(filePath)
	if err == nil {
		// File already exists
		fmt.Printf("File %v already exsits\n", filePath)
		return FileInfo{}, err
	}

	fileInfo := FileInfo{Name: filePath,
		Key: filePath,
		LastModified: time.Now().UTC(),
		UploadID: uuid.New().String(),
		Size: 0}

	err = writeFiletInfo(filePath, &fileInfo)
	return fileInfo, err
}

func WriteFilePart(filePath string, objectPartData []byte, PartNumber int) (FileInfo, error) {

	lock.Lock()
	defer lock.Unlock()

	fileInfo, err := readFileInfo(filePath)
	if err != nil {
		fmt.Printf("Error reading file info: %v\n", err)
		return FileInfo{}, err
	}

	fsutils.AppendToFile(getFilePath(filePath), objectPartData)

	// Update object info
	fileInfo.Size = fileInfo.Size + int64(len(objectPartData))
	fileInfo.LastModified = time.Now().UTC()
	err = writeFiletInfo(filePath, &fileInfo)
	if err != nil {
		fmt.Printf("Error writing file info: %v\n", err)
		return FileInfo{}, err
	}

	return fileInfo, nil
}

func ReadFile(filePath string) ([]byte, error) {

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

func UpdateFileInfo(filePath string, fileInfo FileInfo) error {

	lock.Lock()
	defer lock.Unlock()

	err := writeFiletInfo(filePath, &fileInfo)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(filePath string) error {

	lock.Lock()
	defer lock.Unlock()

	_, err := readFileInfo(filePath)
	if err != nil {
		return err
	}

	_, err = readFileInfo(filePath)
	if err != nil {
		return err
	}

	err = os.Remove(getFileInfoPath(filePath))
	if err != nil {
		return &FileError{Op: "Error deleting object", Key: filePath}
	}

	err = os.Remove(getFilePath(filePath))
	if err != nil {
		return &FileError{Op: "Error deleting object", Key: filePath}
	}

	return nil
}

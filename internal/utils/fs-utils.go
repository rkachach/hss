package fsutils

import (
	"io"
	"time"
	"path/filepath"
	"strings"
	"os"
	"fmt"
)

type EntryInfo struct {
	Name       string
	IsDirectory bool
}

func ListDirectoryEntries(dirPath string) ([]string, error) {
	var entries []string

	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	// Filter directories
	for _, fileInfo := range fileInfos {
		entries = append(entries, fileInfo.Name())
	}

	return entries, nil
}

func ListDirectoryWithDetails(dirPath string) ([]EntryInfo, error) {
	var entries []EntryInfo

	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	// Collect entry information
	for _, fileInfo := range fileInfos {
		entryInfo := EntryInfo{
			Name:       fileInfo.Name(),
			IsDirectory: fileInfo.IsDir(),
		}
		entries = append(entries, entryInfo)
	}

	return entries, nil
}

func ListDirectories(dirPath string) ([]string, error) {
	var directories []string

	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	// Filter directories
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			directories = append(directories, fileInfo.Name())
		}
	}

	return directories, nil
}

func ListFiles(dirPath string) ([]string, error) {
	var files []string

	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	// Filter files
	for _, fileInfo := range fileInfos {
		if ! fileInfo.IsDir() {
			files = append(files, fileInfo.Name())
		}
	}

	return files, nil
}

func DirectoryExists(dirPath string) (bool, error) {
	// Use os.Stat to get information about the file or directory
	// If the directory exists, Stat will return nil (no error)
	// If the directory doesn't exist, it will return an error
	if _, err := os.Stat(dirPath); err == nil {
		// Directory exists
		return true, nil
	} else if os.IsNotExist(err) {
		// Directory does not exist
		return false, nil
	} else {
		// Error occurred while checking directory
		return false, err
	}
}

func AppendToFile(filepath string, data []byte) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to file: %s", err)
	}

	return nil
}

func CopyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    destinationFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destinationFile.Close()

    _, err = io.Copy(destinationFile, sourceFile)
    if err != nil {
        return err
    }

    // Flushes the destination file to ensure all data is written
    err = destinationFile.Sync()
    if err != nil {
        return err
    }

    return nil
}

func SplitPath(input string) (string, string) {
	fmt.Printf("splitting %v \n", input)
	index := strings.Index(input, "/")
	if index == -1 {
		// Handle the case where there's no "/"
		return input, ""
	}

	first := input[:index]
	second := input[index+1:]
	return first, second
}

func GetDirectoryStats(dirPath string) (int64, int, time.Time, error) {
	var size int64
	var filesCount int
	var creationTime time.Time

	err := filepath.Walk(dirPath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			// Skip directories
			return nil
		}

		// Calculate size and count files
		size += fileInfo.Size()
		filesCount++

		// Capture creation time of the directory
		if creationTime.IsZero() {
			creationTime = fileInfo.ModTime()
		}

		return nil
	})

	if err != nil {
		return 0, 0, time.Time{}, err
	}

	return size, filesCount, creationTime, nil
}

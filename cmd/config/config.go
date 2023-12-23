package config

import (
	"fmt"
	"encoding/json"
	"os"
	"log"
	"io"
)

var AppConfig AppConfigRecord
var Logger *log.Logger

type LoggingConfig struct {
	LogFile         string   `json:"log_file"`
	SpecificHeaders []string `json:"specific_headers"`
}

type DataStoreConfig struct {
	Root string   `json:"root"`
}

type AppConfigRecord struct {
	ServerPort   int	      `json:"server_port"`
	ConsolePort  int	      `json:"console_port"`
	Logging      LoggingConfig    `json:"logging"`
	StoreConfig  DataStoreConfig  `json:"object_store"`
}

func ReadConfig(filename string) error {
	fmt.Printf("Loading config from %s\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		return err
	}

	return nil
}

func InitLogger() {

	// Open a file for logging
	logFile, err := os.OpenFile(AppConfig.Logging.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Use log.New to create a new logger for writing logs to the file
	Logger = log.New(logFile, "", log.Ldate|log.Ltime)
	// Create a multi-writer to write logs to both stdout and the file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	// Create a logger to write logs to the multi-writer
	Logger = log.New(multiWriter, "", log.Ldate|log.Ltime)
}

package utils

import (
	"log"
	"os"
)

var logFile *os.File

func InitLogger(logFilePath string) {
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

func LogError(format string, v ...interface{}) {
	log.Printf("ERROR: "+format, v...)
}

func LogInfo(format string, v ...interface{}) {
	log.Printf("INFO: "+format, v...)
}

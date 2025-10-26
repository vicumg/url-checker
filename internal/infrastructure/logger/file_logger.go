package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FileLogger struct {
	file *os.File
}

func NewFileLogger(logDir string) (*FileLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("checks_%s.log", time.Now().Format("2006-01-02")))

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FileLogger{file: file}, nil
}

func (l *FileLogger) LogCheck(monitorID, url string, statusCode int, responseTime time.Duration, err error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if err != nil {
		log.Printf("[%s] [%s] %s - ERROR: %v\n", timestamp, monitorID, url, err)
		fmt.Fprintf(l.file, "[%s] [%s] %s - ERROR: %v\n", timestamp, monitorID, url, err)
	} else {
		log.Printf("[%s] [%s] %s - Status: %d, Time: %v\n", timestamp, monitorID, url, statusCode, responseTime)
		fmt.Fprintf(l.file, "[%s] [%s] %s - Status: %d, Time: %v\n", timestamp, monitorID, url, statusCode, responseTime)
	}
}

func (l *FileLogger) Close() error {
	return l.file.Close()
}

package logging

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// generateRandomString creates a random alphanumeric string of length 6
func generateRandomString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// CreateLog writes a log entry to a daily log file
func MakeLog(moduleName, request, response interface{}) error {
	today := time.Now().Format("020106") // Format: ddMMyy
	dir := "logs"                        // Directory to store logs
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Check if a file for today already exists
	logFilePattern := fmt.Sprintf("log-*-%s.txt", today)
	files, err := filepath.Glob(filepath.Join(dir, logFilePattern))
	if err != nil {
		return fmt.Errorf("failed to search log files: %v", err)
	}

	var logFileName string
	if len(files) == 0 {
		// Create a new log file with a random string
		logFileName = filepath.Join(dir, fmt.Sprintf("log-%s-%s.txt", generateRandomString(6), today))
	} else {
		// Use the existing log file
		logFileName = files[0]
	}

	// Open the file in append mode, create if not exists
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	// Write log message with moduleName, timestamp, request, and response
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] Module: %s | Request: %s | Response: %s\n", timestamp, moduleName, request, response)
	_, err = file.WriteString(logEntry)
	if err != nil {
		return fmt.Errorf("failed to write log: %v", err)
	}
	fmt.Println(logEntry)

	return nil
}

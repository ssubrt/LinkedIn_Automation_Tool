package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// Info logs an informational message
func Info(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s[%s] %sINFO%s: %s\n", colorBlue, timestamp, colorGreen, colorReset, message)
}

// Error logs an error message
func Error(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s[%s] %sERROR%s: %s\n", colorRed, timestamp, colorRed, colorReset, message)
	log.Printf("[%s] ERROR: %s", timestamp, message)
}

// Warning logs a warning message
func Warning(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s[%s] %sWARNING%s: %s\n", colorYellow, timestamp, colorYellow, colorReset, message)
}

// Debug logs a debug message
func Debug(message string) {
	if os.Getenv("DEBUG") == "true" {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("%s[%s] %sDEBUG%s: %s\n", colorBlue, timestamp, colorBlue, colorReset, message)
	}
}

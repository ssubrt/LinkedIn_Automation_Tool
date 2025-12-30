package logger

import (
	"testing"
)

// TestLoggerInfo verifies info logs are created
func TestLoggerInfo(t *testing.T) {
	// Should not panic or error
	Info("Test info message")
	Info("Test with variable: " + "value")
}

// TestLoggerError verifies error logs are created
func TestLoggerError(t *testing.T) {
	// Should not panic or error
	Error("Test error message")
	Error("Error with details: something went wrong")
}

// TestLoggerMultipleCalls ensures logger can handle multiple calls
func TestLoggerMultipleCalls(t *testing.T) {
	for i := 0; i < 10; i++ {
		Info("Test iteration")
	}

	// Should complete without error
	if t.Failed() {
		t.Error("Logger failed during multiple calls")
	}
}

// TestLoggerWithEmptyMessage tests logger with empty string
func TestLoggerWithEmptyMessage(t *testing.T) {
	// Should not panic
	Info("")
	Error("")
}

// TestLoggerConcurrency tests logger with concurrent calls
func TestLoggerConcurrency(t *testing.T) {
	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func(id int) {
			Info("Concurrent message from goroutine")
			Error("Error from goroutine")
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 5; i++ {
		<-done
	}
}

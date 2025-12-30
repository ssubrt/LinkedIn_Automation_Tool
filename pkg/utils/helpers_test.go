package utils

import (
	"testing"
	"time"
)

// TestGenerateRandomDelay tests random delay generation
func TestGenerateRandomDelay(t *testing.T) {
	delay := GenerateRandomDelay(100, 200)

	// Should return a duration
	if delay < 100*time.Millisecond || delay > 200*time.Millisecond {
		t.Errorf("GenerateRandomDelay(100, 200) = %v, want between 100ms and 200ms", delay)
	}
}

// TestGenerateRandomCoordinates tests random coordinate generation
func TestGenerateRandomCoordinates(t *testing.T) {
	x, y := GenerateRandomCoordinates(100, 800, 100, 600)

	if x < 100 || x > 800 {
		t.Errorf("GenerateRandomCoordinates X = %d, want between 100 and 800", x)
	}

	if y < 100 || y > 600 {
		t.Errorf("GenerateRandomCoordinates Y = %d, want between 100 and 600", y)
	}
}

// TestGenerateRandomScrollDistance tests scroll distance generation
func TestGenerateRandomScrollDistance(t *testing.T) {
	dist := GenerateRandomScrollDistance(200, 600)

	if dist < 200 || dist > 600 {
		t.Errorf("GenerateRandomScrollDistance = %d, want between 200 and 600", dist)
	}
}

// TestGenerateSessionID tests session ID generation
func TestGenerateSessionID(t *testing.T) {
	id := GenerateSessionID()

	if id == "" {
		t.Error("GenerateSessionID() returned empty string")
	}

	if len(id) < 10 {
		t.Errorf("GenerateSessionID() returned %q, expected longer ID", id)
	}
}

// TestFormatDuration tests duration formatting
func TestFormatDuration(t *testing.T) {
	tests := []struct {
		ms       int64
		expected string
	}{
		{500, "500ms"},
		{1000, "1s"},
		{60000, "1m 0s"},
		{65000, "1m 5s"},
		{3600000, "60m 0s"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := FormatDuration(tt.ms)
			if got != tt.expected {
				t.Errorf("FormatDuration(%d) = %q, want %q", tt.ms, got, tt.expected)
			}
		})
	}
}

// TestContainsString tests string in slice search
func TestContainsString(t *testing.T) {
	tests := []struct {
		slice    []string
		str      string
		expected bool
	}{
		{[]string{"a", "b", "c"}, "b", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{}, "a", false},
		{[]string{"a"}, "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			got := ContainsString(tt.slice, tt.str)
			if got != tt.expected {
				t.Errorf("ContainsString(%v, %q) = %v, want %v", tt.slice, tt.str, got, tt.expected)
			}
		})
	}
}

// TestReverseString tests string reversal
func TestReverseString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "olleh"},
		{"a", "a"},
		{"abc", "cba"},
		{"", ""},
		{"12345", "54321"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ReverseString(tt.input)
			if got != tt.expected {
				t.Errorf("ReverseString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateRandomDelay creates a random delay within range
func GenerateRandomDelay(minMs, maxMs int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	delay := r.Intn(maxMs-minMs+1) + minMs
	return time.Duration(delay) * time.Millisecond
}

// GenerateRandomCoordinates creates random X, Y coordinates
func GenerateRandomCoordinates(minX, maxX, minY, maxY int) (int, int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	x := r.Intn(maxX-minX+1) + minX
	y := r.Intn(maxY-minY+1) + minY
	return x, y
}

// GenerateRandomScrollDistance creates random scroll distance
func GenerateRandomScrollDistance(minDist, maxDist int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(maxDist-minDist+1) + minDist
}

// GenerateSessionID creates a unique session identifier
func GenerateSessionID() string {
	return fmt.Sprintf("session_%d_%d", time.Now().Unix(), rand.Intn(10000))
}

// FormatDuration formats milliseconds to human-readable string
func FormatDuration(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}
	seconds := ms / 1000
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	minutes := seconds / 60
	remainingSeconds := seconds % 60
	return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
}

// ContainsString checks if string exists in slice
func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// ReverseString reverses a string
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsLinkedInCheckpoint checks if the current URL is a LinkedIn verification/checkpoint page
// These pages appear when LinkedIn suspects automation and requires manual verification
func IsLinkedInCheckpoint(url string) bool {
	checkpointPatterns := []string{
		"/checkpoint/",
		"/challenge/",
		"/uas/login-verification",
		"/uas/challenge",
		"/cap/", // CAPTCHA page
	}

	for _, pattern := range checkpointPatterns {
		if len(url) > 0 && ContainsString([]string{url}, pattern) {
			return true
		}
	}
	return false
}

// ExtractProfileID extracts the profile ID from a LinkedIn URL
func ExtractProfileID(url string) string {
	// Remove query parameters
	// URLs are typically https://www.linkedin.com/in/profile-id/
	// or /in/profile-id/

	// Find /in/
	inIdx := -1
	for i := 0; i < len(url)-3; i++ {
		if url[i:i+4] == "/in/" {
			inIdx = i
			break
		}
	}

	if inIdx != -1 {
		start := inIdx + 4
		end := len(url)

		// Find next slash or ?
		for i := start; i < len(url); i++ {
			if url[i] == '/' || url[i] == '?' {
				end = i
				break
			}
		}

		if start < end {
			return url[start:end]
		}
	}

	return ""
}

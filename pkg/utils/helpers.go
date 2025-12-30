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

package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}

	// Basic email regex validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword checks if password meets requirements
func ValidatePassword(password string) bool {
	// Minimum 6 characters
	if len(password) < 6 {
		return false
	}

	// Maximum 128 characters
	if len(password) > 128 {
		return false
	}

	return true
}

// ValidateURL checks if URL format is valid
func ValidateURL(url string) bool {
	urlRegex := regexp.MustCompile(`^https?://`)
	return urlRegex.MatchString(url)
}

// IsLinkedInURL checks if URL is from LinkedIn
func IsLinkedInURL(url string) bool {
	return strings.Contains(url, "linkedin.com")
}

// SanitizeEmail removes whitespace and converts to lowercase
func SanitizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// SanitizePassword removes leading/trailing whitespace only
func SanitizePassword(password string) string {
	return strings.TrimSpace(password)
}

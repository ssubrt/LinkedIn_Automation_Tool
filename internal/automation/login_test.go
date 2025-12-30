package automation

import (
	"testing"
)

// TestCredentialValidation verifies credentials format
func TestCredentialValidation(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		valid    bool
	}{
		{"Valid email and password", "user@example.com", "password123", true},
		{"Valid Gmail", "user@gmail.com", "Password@123", true},
		{"Valid LinkedIn email", "user@linkedin.com", "SecurePass123", true},
		{"Empty email", "", "password123", false},
		{"Empty password", "user@example.com", "", false},
		{"Both empty", "", "", false},
		{"Very long password", "user@example.com", "VeryLongPasswordWith123456789SpecialChars!@#$%^&*()", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate email format (basic check)
			valid := isValidEmail(tt.email) && tt.password != ""
			if valid != tt.valid {
				t.Errorf("Expected %v, got %v for email=%s, password=%s", tt.valid, valid, tt.email, tt.password)
			}
		})
	}
}

// TestEmailValidation tests email format validation
func TestEmailValidation(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"user@example.com", true},
		{"test@gmail.com", true},
		{"test@linkedin.com", true},
		{"invalid.email", false},
		{"@example.com", false},
		{"user@", false},
		{"", false},
		{"user @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			valid := isValidEmail(tt.email)
			if valid != tt.valid {
				t.Errorf("Expected %v, got %v for email=%s", tt.valid, valid, tt.email)
			}
		})
	}
}

// TestPasswordValidation tests password requirements
func TestPasswordValidation(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"password123", true},
		{"SecurePass@123", true},
		{"VeryLongPassword1", true},
		{"", false},
		{"short", false},
		{"123456", true},
		{"pass", false},
		{"12345", false},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			valid := isValidPassword(tt.password)
			if valid != tt.valid {
				t.Errorf("Expected %v, got %v for password=%s", tt.valid, valid, tt.password)
			}
		})
	}
}

// Helper function for email validation - stricter version
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}

	// Must have minimum length
	if len(email) < 5 {
		return false
	}

	// Count @ symbols (must be exactly 1)
	atCount := 0
	atIndex := -1
	for i, ch := range email {
		if ch == '@' {
			atCount++
			atIndex = i
		}
	}

	if atCount != 1 {
		return false
	}

	// Must have content before @
	if atIndex == 0 {
		return false
	}

	// Must have content after @
	if atIndex == len(email)-1 {
		return false
	}

	// Username part (before @) - no spaces allowed
	username := email[:atIndex]
	for _, ch := range username {
		if ch == ' ' {
			return false
		}
	}

	// Domain part (after @) - must end with valid TLD
	domain := email[atIndex+1:]
	validTLDs := []string{".com", ".org", ".in", ".net", ".co.uk", ".edu"}
	hasValidTLD := false
	for _, tld := range validTLDs {
		if len(domain) >= len(tld) && domain[len(domain)-len(tld):] == tld {
			hasValidTLD = true
			break
		}
	}

	return hasValidTLD
}

// Helper function for password validation
func isValidPassword(password string) bool {
	// Minimum length of 6 characters
	return len(password) >= 6
}

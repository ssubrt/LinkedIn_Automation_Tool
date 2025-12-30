package utils

import (
	"testing"
)

// TestValidateEmail tests email validation
func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"user@example.com", true},
		{"test.user@gmail.com", true},
		{"user+tag@linkedin.com", true},
		{"", false},
		{"invalid", false},
		{"user@", false},
		{"@example.com", false},
		{"user @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			got := ValidateEmail(tt.email)
			if got != tt.valid {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.valid)
			}
		})
	}
}

// TestValidatePassword tests password validation
func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"password123", true},
		{"SecurePass@123", true},
		{"VeryLongPassword1", true},
		{"", false},
		{"short", false},
		{"12345", false},
		{"123456", true},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			got := ValidatePassword(tt.password)
			if got != tt.valid {
				t.Errorf("ValidatePassword(%q) = %v, want %v", tt.password, got, tt.valid)
			}
		})
	}
}

// TestValidateURL tests URL validation
func TestValidateURL(t *testing.T) {
	tests := []struct {
		url   string
		valid bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"ftp://example.com", false},
		{"example.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := ValidateURL(tt.url)
			if got != tt.valid {
				t.Errorf("ValidateURL(%q) = %v, want %v", tt.url, got, tt.valid)
			}
		})
	}
}

// TestIsLinkedInURL tests LinkedIn URL detection
func TestIsLinkedInURL(t *testing.T) {
	tests := []struct {
		url      string
		isLinked bool
	}{
		{"https://linkedin.com/login", true},
		{"https://www.linkedin.com/feed", true},
		{"https://google.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := IsLinkedInURL(tt.url)
			if got != tt.isLinked {
				t.Errorf("IsLinkedInURL(%q) = %v, want %v", tt.url, got, tt.isLinked)
			}
		})
	}
}

// TestSanitizeEmail tests email sanitization
func TestSanitizeEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"User@EXAMPLE.COM", "user@example.com"},
		{"  user@example.com  ", "user@example.com"},
		{"USER@EXAMPLE.COM", "user@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := SanitizeEmail(tt.input)
			if got != tt.expected {
				t.Errorf("SanitizeEmail(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestSanitizePassword tests password sanitization
func TestSanitizePassword(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  password123  ", "password123"},
		{"password123", "password123"},
		{"\tpassword123\n", "password123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := SanitizePassword(tt.input)
			if got != tt.expected {
				t.Errorf("SanitizePassword(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

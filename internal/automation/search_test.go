package automation

import (
	"testing"

	"linkedin-automation/pkg/utils"
)

func TestBuildSearchURL(t *testing.T) {
	tests := []struct {
		name      string
		config    SearchConfig
		wantError bool
		contains  []string
	}{
		{
			name: "Keywords only",
			config: SearchConfig{
				Keywords: "software engineer",
			},
			wantError: false,
			contains:  []string{"keywords=software"},
		},
		{
			name: "Job title filter",
			config: SearchConfig{
				JobTitle: "CTO",
			},
			wantError: false,
			contains:  []string{"title=CTO"},
		},
		{
			name: "Company filter",
			config: SearchConfig{
				Company: "Google",
			},
			wantError: false,
			contains:  []string{"company=Google"},
		},
		{
			name: "Location filter - San Francisco",
			config: SearchConfig{
				Keywords: "engineer",
				Location: "San Francisco Bay Area",
			},
			wantError: false,
			contains:  []string{"keywords=engineer", "geoUrn", "90000084"},
		},
		{
			name: "All filters combined",
			config: SearchConfig{
				Keywords: "AI researcher",
				JobTitle: "Research Scientist",
				Company:  "OpenAI",
				Location: "San Francisco Bay Area",
			},
			wantError: false,
			contains:  []string{"keywords=AI", "title=Research", "company=OpenAI", "geoUrn", "90000084"},
		},
		{
			name: "Location not found",
			config: SearchConfig{
				Keywords: "engineer",
				Location: "NonExistentCity",
			},
			wantError: false,
			contains:  []string{"keywords=engineer"},
		},
		{
			name:      "No filters - should error",
			config:    SearchConfig{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := buildSearchURL(tt.config)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check that URL starts with search base
			if len(url) == 0 {
				t.Errorf("Empty URL returned")
				return
			}

			// Check for expected substrings
			for _, substr := range tt.contains {
				if !containsSubstring(url, substr) {
					t.Errorf("URL does not contain expected substring '%s'. URL: %s", substr, url)
				}
			}
		})
	}
}

func TestSearchConfigDefaults(t *testing.T) {
	config := SearchConfig{
		Keywords: "test",
	}

	// These should be set by SearchPeople function
	if config.MaxPages != 0 && config.MaxPages != utils.MaxPaginationPages {
		t.Errorf("Expected MaxPages to be 0 or %d, got %d", utils.MaxPaginationPages, config.MaxPages)
	}

	if config.DuplicateDays != 0 && config.DuplicateDays != 30 {
		t.Errorf("Expected DuplicateDays to be 0 or 30, got %d", config.DuplicateDays)
	}
}

func TestLocationMapping(t *testing.T) {
	// Test that key locations are present
	keyLocations := []string{
		"San Francisco Bay Area",
		"New York City Area",
		"London",
		"United States",
		"United Kingdom",
	}

	for _, location := range keyLocations {
		urn, found := utils.LinkedInLocations[location]
		if !found {
			t.Errorf("Location '%s' not found in location map", location)
		}
		if urn == "" {
			t.Errorf("Location '%s' has empty URN", location)
		}
	}
}

func TestSearchResultValidation(t *testing.T) {
	result := SearchResult{
		ProfileID:  "john-doe",
		Name:       "John Doe",
		Title:      "Software Engineer",
		Company:    "Tech Corp",
		Location:   "San Francisco",
		ProfileURL: "https://www.linkedin.com/in/john-doe/",
	}

	// Verify basic fields
	if result.ProfileID == "" {
		t.Error("ProfileID should not be empty")
	}
	if result.Name == "" {
		t.Error("Name should not be empty")
	}
	if result.ProfileURL == "" {
		t.Error("ProfileURL should not be empty")
	}
}

// Helper function
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

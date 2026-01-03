package automation

import (
	"strings"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    MessageTemplate
		vars        TemplateVariables
		wantError   bool
		contains    []string
		notContains []string
	}{
		{
			name: "Connection request - generic",
			template: MessageTemplate{
				ID:        "test_conn",
				Type:      TemplateConnectionRequest,
				Body:      "Hi {{.FirstName}}, I saw you work at {{.Company}}. Let's connect!",
				MaxLength: ConnectionNoteMaxLength,
			},
			vars: TemplateVariables{
				FirstName: "John",
				Company:   "Google",
			},
			wantError:   false,
			contains:    []string{"Hi John", "work at Google"},
			notContains: []string{"{{.FirstName}}", "{{.Company}}"},
		},
		{
			name: "Connection request - with sender info",
			template: MessageTemplate{
				ID:        "test_conn2",
				Type:      TemplateConnectionRequest,
				Body:      "Hi {{.FirstName}}, I'm {{.YourName}} from {{.YourCompany}}. Let's connect!",
				MaxLength: ConnectionNoteMaxLength,
			},
			vars: TemplateVariables{
				FirstName:   "Jane",
				YourName:    "Bob Smith",
				YourCompany: "Microsoft",
			},
			wantError: false,
			contains:  []string{"Hi Jane", "I'm Bob Smith", "from Microsoft"},
		},
		{
			name: "Message - introduction",
			template: MessageTemplate{
				ID:        "test_msg",
				Type:      TemplateIntroduction,
				Body:      "Hi {{.FirstName}},\n\nI'm {{.YourName}}, {{.YourTitle}} at {{.YourCompany}}.\n\nBest regards,\n{{.YourName}}",
				MaxLength: MessageMaxLength,
			},
			vars: TemplateVariables{
				FirstName:   "Sarah",
				YourName:    "Mike Johnson",
				YourTitle:   "CTO",
				YourCompany: "TechCorp",
			},
			wantError: false,
			contains:  []string{"Hi Sarah", "I'm Mike Johnson", "CTO at TechCorp", "Best regards"},
		},
		{
			name: "Message too long",
			template: MessageTemplate{
				ID:        "test_long",
				Type:      TemplateConnectionRequest,
				Body:      strings.Repeat("A", 350), // Exceeds 300 char limit
				MaxLength: ConnectionNoteMaxLength,
			},
			vars:      TemplateVariables{},
			wantError: true,
		},
		{
			name: "Empty template",
			template: MessageTemplate{
				ID:        "test_empty",
				Type:      TemplateConnectionRequest,
				Body:      "{{.CustomReason}}", // Will be empty
				MaxLength: ConnectionNoteMaxLength,
			},
			vars:      TemplateVariables{}, // No CustomReason provided
			wantError: true,                // Should error on empty message
		},
		{
			name: "Auto-extract first name",
			template: MessageTemplate{
				ID:        "test_auto",
				Type:      TemplateConnectionRequest,
				Body:      "Hi {{.FirstName}}!",
				MaxLength: ConnectionNoteMaxLength,
			},
			vars: TemplateVariables{
				FullName: "John Doe",
			},
			wantError: false,
			contains:  []string{"Hi John!"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderTemplate(tt.template, tt.vars)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none. Result: %s", result)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check for expected strings
			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("Result does not contain expected string '%s'. Result: %s", s, result)
				}
			}

			// Check for strings that should not be present
			for _, s := range tt.notContains {
				if strings.Contains(result, s) {
					t.Errorf("Result contains unexpected string '%s'. Result: %s", s, result)
				}
			}

			// Verify length is within limits
			if len(result) > tt.template.MaxLength {
				t.Errorf("Result exceeds max length: %d > %d", len(result), tt.template.MaxLength)
			}
		})
	}
}

func TestRenderSubject(t *testing.T) {
	tests := []struct {
		name     string
		subject  string
		vars     TemplateVariables
		expected string
	}{
		{
			name:    "Simple subject",
			subject: "Hello {{.FirstName}}!",
			vars: TemplateVariables{
				FirstName: "Alice",
			},
			expected: "Hello Alice!",
		},
		{
			name:    "Subject with company",
			subject: "Opportunity at {{.Company}}",
			vars: TemplateVariables{
				Company: "Amazon",
			},
			expected: "Opportunity at Amazon",
		},
		{
			name:     "Subject too long gets truncated",
			subject:  strings.Repeat("A", 250),
			vars:     TemplateVariables{},
			expected: strings.Repeat("A", SubjectMaxLength-3) + "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderSubject(tt.subject, tt.vars)

			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}

			if len(result) > SubjectMaxLength {
				t.Errorf("Subject exceeds max length: %d > %d", len(result), SubjectMaxLength)
			}
		})
	}
}

func TestGetTemplateByID(t *testing.T) {
	tests := []struct {
		name       string
		templateID string
		wantError  bool
		wantType   TemplateType
	}{
		{
			name:       "Valid connection template",
			templateID: "conn_generic",
			wantError:  false,
			wantType:   TemplateConnectionRequest,
		},
		{
			name:       "Valid message template",
			templateID: "msg_introduction",
			wantError:  false,
			wantType:   TemplateIntroduction,
		},
		{
			name:       "Invalid template ID",
			templateID: "nonexistent",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := GetTemplateByID(tt.templateID)

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

			if template.Type != tt.wantType {
				t.Errorf("Expected type %s, got %s", tt.wantType, template.Type)
			}

			if template.ID != tt.templateID {
				t.Errorf("Expected ID %s, got %s", tt.templateID, template.ID)
			}
		})
	}
}

func TestGetTemplatesByType(t *testing.T) {
	connectionTemplates := GetTemplatesByType(TemplateConnectionRequest)
	if len(connectionTemplates) == 0 {
		t.Error("Expected at least one connection template")
	}

	for _, template := range connectionTemplates {
		if template.Type != TemplateConnectionRequest {
			t.Errorf("Expected connection template, got %s", template.Type)
		}
	}

	messageTemplates := GetTemplatesByType(TemplateIntroduction)
	for _, template := range messageTemplates {
		if template.Type != TemplateIntroduction {
			t.Errorf("Expected introduction template, got %s", template.Type)
		}
	}
}

func TestValidateMessageLength(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		messageType TemplateType
		wantError   bool
	}{
		{
			name:        "Valid connection note",
			message:     "Hi, let's connect!",
			messageType: TemplateConnectionRequest,
			wantError:   false,
		},
		{
			name:        "Connection note too long",
			message:     strings.Repeat("A", 350),
			messageType: TemplateConnectionRequest,
			wantError:   true,
		},
		{
			name:        "Valid message",
			message:     strings.Repeat("A", 100),
			messageType: TemplateIntroduction,
			wantError:   false,
		},
		{
			name:        "Message too long",
			message:     strings.Repeat("A", 9000),
			messageType: TemplateIntroduction,
			wantError:   true,
		},
		{
			name:        "Empty message",
			message:     "",
			messageType: TemplateConnectionRequest,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMessageLength(tt.message, tt.messageType)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCleanupWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Multiple spaces",
			input:    "Hello    world",
			expected: "Hello world",
		},
		{
			name:     "Multiple newlines",
			input:    "Line 1\n\n\n\nLine 2",
			expected: "Line 1\n\nLine 2",
		},
		{
			name:     "Leading/trailing whitespace",
			input:    "  Hello  \n  World  ",
			expected: "Hello\nWorld",
		},
		{
			name:     "Normal text",
			input:    "Hello world",
			expected: "Hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanupWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestTruncateMessage(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		maxLength int
		expected  string
	}{
		{
			name:      "No truncation needed",
			message:   "Short message",
			maxLength: 100,
			expected:  "Short message",
		},
		{
			name:      "Truncation needed",
			message:   "This is a very long message that needs to be truncated",
			maxLength: 20,
			expected:  "This is a very lo...",
		},
		{
			name:      "Exact length",
			message:   "12345",
			maxLength: 5,
			expected:  "12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateMessage(tt.message, tt.maxLength)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
			if len(result) > tt.maxLength {
				t.Errorf("Result exceeds max length: %d > %d", len(result), tt.maxLength)
			}
		})
	}
}

func TestGetConnectionRequestTemplates(t *testing.T) {
	templates := GetConnectionRequestTemplates()

	if len(templates) == 0 {
		t.Error("Expected at least one connection template")
	}

	for _, template := range templates {
		// Verify required fields
		if template.ID == "" {
			t.Error("Template ID is empty")
		}
		if template.Body == "" {
			t.Error("Template body is empty")
		}
		if template.Type != TemplateConnectionRequest {
			t.Errorf("Expected connection request type, got %s", template.Type)
		}
		if template.MaxLength != ConnectionNoteMaxLength {
			t.Errorf("Expected max length %d, got %d", ConnectionNoteMaxLength, template.MaxLength)
		}
	}
}

func TestGetMessageTemplates(t *testing.T) {
	templates := GetMessageTemplates()

	if len(templates) == 0 {
		t.Error("Expected at least one message template")
	}

	for _, template := range templates {
		// Verify required fields
		if template.ID == "" {
			t.Error("Template ID is empty")
		}
		if template.Body == "" {
			t.Error("Template body is empty")
		}
		if template.Type == TemplateConnectionRequest {
			t.Error("Message template should not have connection request type")
		}
		if template.MaxLength != MessageMaxLength {
			t.Errorf("Expected max length %d, got %d", MessageMaxLength, template.MaxLength)
		}
	}
}

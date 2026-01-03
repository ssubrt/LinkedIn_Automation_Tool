package automation

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"linkedin-automation/internal/logger"
)

// TemplateType represents the type of message template
type TemplateType string

const (
	TemplateConnectionRequest TemplateType = "connection_request"
	TemplateFollowUp          TemplateType = "follow_up"
	TemplateIntroduction      TemplateType = "introduction"
	TemplateNetworking        TemplateType = "networking"
)

// TemplateVariables holds variables for template substitution
type TemplateVariables struct {
	FirstName    string // Recipient's first name
	LastName     string // Recipient's last name
	FullName     string // Recipient's full name
	Title        string // Recipient's job title
	Company      string // Recipient's company
	Industry     string // Industry/sector
	YourName     string // Sender's name
	YourTitle    string // Sender's title
	YourCompany  string // Sender's company
	CustomReason string // Custom reason for connection
	Date         string // Current date
}

// MessageTemplate represents a message template with metadata
type MessageTemplate struct {
	ID          string
	Type        TemplateType
	Name        string
	Subject     string // For messages only (not used in connection requests)
	Body        string
	Description string
	MaxLength   int // Character limit (300 for connection notes, 8000 for messages)
}

// Character limits per LinkedIn's specifications
const (
	ConnectionNoteMaxLength = 300  // LinkedIn's limit for connection request notes
	MessageMaxLength        = 8000 // LinkedIn's limit for direct messages
	SubjectMaxLength        = 200  // LinkedIn's limit for message subjects
)

// GetConnectionRequestTemplates returns predefined connection request templates
func GetConnectionRequestTemplates() []MessageTemplate {
	return []MessageTemplate{
		{
			ID:          "conn_generic",
			Type:        TemplateConnectionRequest,
			Name:        "Generic Professional",
			Body:        "Hi {{.FirstName}}, I came across your profile and was impressed by your work at {{.Company}}. I'd love to connect{{if .Industry}} and learn more about your experience in {{.Industry}}{{end}}.",
			Description: "Generic professional connection request",
			MaxLength:   ConnectionNoteMaxLength,
		},
		{
			ID:          "conn_role_specific",
			Type:        TemplateConnectionRequest,
			Name:        "Role-Specific",
			Body:        "Hi {{.FirstName}}, I noticed you're a {{.Title}} at {{.Company}}. I'm {{.YourTitle}} at {{.YourCompany}} and would love to connect to exchange insights about our field.",
			Description: "Connection based on similar roles",
			MaxLength:   ConnectionNoteMaxLength,
		},
		{
			ID:          "conn_industry",
			Type:        TemplateConnectionRequest,
			Name:        "Industry Connection",
			Body:        "Hi {{.FirstName}}, I saw your profile and noticed we both work in {{.Industry}}. I'd appreciate the opportunity to connect and potentially collaborate in the future.",
			Description: "Connection based on shared industry",
			MaxLength:   ConnectionNoteMaxLength,
		},
		{
			ID:          "conn_mutual_interest",
			Type:        TemplateConnectionRequest,
			Name:        "Mutual Interest",
			Body:        "Hi {{.FirstName}}, your experience at {{.Company}} caught my attention. {{.CustomReason}} I'd love to connect and learn from your expertise.",
			Description: "Connection with custom reason",
			MaxLength:   ConnectionNoteMaxLength,
		},
		{
			ID:          "conn_networking",
			Type:        TemplateConnectionRequest,
			Name:        "Networking",
			Body:        "Hi {{.FirstName}}, I'm expanding my professional network with {{.Industry}} professionals. Your background at {{.Company}} is impressive. Let's connect!",
			Description: "General networking connection",
			MaxLength:   ConnectionNoteMaxLength,
		},
		{
			ID:          "conn_brief",
			Type:        TemplateConnectionRequest,
			Name:        "Brief & Direct",
			Body:        "Hi {{.FirstName}}, impressive work at {{.Company}}! Would love to connect.",
			Description: "Short and direct connection request",
			MaxLength:   ConnectionNoteMaxLength,
		},
	}
}

// GetMessageTemplates returns predefined message templates
func GetMessageTemplates() []MessageTemplate {
	return []MessageTemplate{
		{
			ID:          "msg_introduction",
			Type:        TemplateIntroduction,
			Name:        "Professional Introduction",
			Subject:     "Great to connect, {{.FirstName}}!",
			Body:        "Hi {{.FirstName}},\n\nThank you for connecting! I'm {{.YourName}}, {{.YourTitle}} at {{.YourCompany}}.\n\nI was impressed by your work as {{.Title}} at {{.Company}}. I'd love to learn more about your experience and explore potential collaboration opportunities.\n\nLooking forward to staying in touch!\n\nBest regards,\n{{.YourName}}",
			Description: "Initial message after connection",
			MaxLength:   MessageMaxLength,
		},
		{
			ID:          "msg_follow_up",
			Type:        TemplateFollowUp,
			Name:        "Follow-Up Message",
			Subject:     "Following up on my previous message",
			Body:        "Hi {{.FirstName}},\n\nI wanted to follow up on my previous message. I'm still very interested in learning about your experience at {{.Company}}.\n\n{{.CustomReason}}\n\nWould you be open to a brief conversation?\n\nBest regards,\n{{.YourName}}",
			Description: "Follow-up after no response",
			MaxLength:   MessageMaxLength,
		},
		{
			ID:          "msg_networking",
			Type:        TemplateNetworking,
			Name:        "Networking Opportunity",
			Subject:     "Exploring opportunities in {{.Industry}}",
			Body:        "Hi {{.FirstName}},\n\nI hope this message finds you well. I'm reaching out to professionals in {{.Industry}} to expand my network and learn from experienced leaders like yourself.\n\nYour background as {{.Title}} at {{.Company}} is particularly interesting to me. Would you be open to sharing some insights about your career journey?\n\nI'd be happy to schedule a brief call at your convenience.\n\nThank you for your time!\n\nBest regards,\n{{.YourName}}\n{{.YourTitle}} at {{.YourCompany}}",
			Description: "Networking and career advice",
			MaxLength:   MessageMaxLength,
		},
		{
			ID:          "msg_collaboration",
			Type:        TemplateNetworking,
			Name:        "Collaboration Proposal",
			Subject:     "Potential collaboration opportunity",
			Body:        "Hi {{.FirstName}},\n\nI came across your profile and was impressed by your work at {{.Company}}.\n\n{{.CustomReason}}\n\nI believe there might be synergies between what you're doing and my work at {{.YourCompany}}. Would you be interested in exploring potential collaboration opportunities?\n\nI'd love to schedule a brief call to discuss further.\n\nLooking forward to hearing from you!\n\nBest regards,\n{{.YourName}}",
			Description: "Business collaboration proposal",
			MaxLength:   MessageMaxLength,
		},
		{
			ID:          "msg_value_add",
			Type:        TemplateIntroduction,
			Name:        "Value-Add Introduction",
			Subject:     "Quick introduction from {{.YourName}}",
			Body:        "Hi {{.FirstName}},\n\nI recently came across your profile and thought you might be interested in {{.CustomReason}}.\n\nAs {{.YourTitle}} at {{.YourCompany}}, I've been working on similar challenges and would be happy to share some insights that might be helpful.\n\nWould you be open to a quick chat?\n\nBest regards,\n{{.YourName}}",
			Description: "Offering value or insights",
			MaxLength:   MessageMaxLength,
		},
	}
}

// RenderTemplate renders a template with the given variables
func RenderTemplate(tmplDef MessageTemplate, vars TemplateVariables) (string, error) {
	// Set default values if not provided
	if vars.FullName == "" && vars.FirstName != "" {
		if vars.LastName != "" {
			vars.FullName = vars.FirstName + " " + vars.LastName
		} else {
			vars.FullName = vars.FirstName
		}
	}

	// Set current date if not provided
	if vars.Date == "" {
		vars.Date = time.Now().Format("January 2, 2006")
	}

	// Extract first name if not provided
	if vars.FirstName == "" && vars.FullName != "" {
		parts := strings.Split(vars.FullName, " ")
		vars.FirstName = parts[0]
		if len(parts) > 1 {
			vars.LastName = strings.Join(parts[1:], " ")
		}
	}

	// Parse the template
	t, err := template.New(tmplDef.ID).Parse(tmplDef.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	result := buf.String()

	// Clean up extra whitespace
	result = cleanupWhitespace(result)

	// Validate length
	if len(result) > tmplDef.MaxLength {
		return "", fmt.Errorf("rendered message exceeds maximum length (%d > %d)", len(result), tmplDef.MaxLength)
	}

	// Validate that we didn't end up with an empty message
	if strings.TrimSpace(result) == "" {
		return "", fmt.Errorf("rendered message is empty - check that template variables are provided")
	}

	logger.Info(fmt.Sprintf("Rendered template '%s' (%d characters)", tmplDef.Name, len(result)))
	return result, nil
}

// RenderSubject renders a message subject line with variables
func RenderSubject(subjectTemplate string, vars TemplateVariables) string {
	// Set defaults
	if vars.FullName == "" && vars.FirstName != "" {
		vars.FullName = vars.FirstName
	}
	if vars.FirstName == "" && vars.FullName != "" {
		parts := strings.Split(vars.FullName, " ")
		vars.FirstName = parts[0]
	}

	// Parse the template
	t, err := template.New("subject").Parse(subjectTemplate)
	if err != nil {
		// Fallback to simple replacement if parsing fails
		logger.Warning("Failed to parse subject template, falling back to simple replacement: " + err.Error())
		return subjectTemplate
	}

	// Execute the template
	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		logger.Warning("Failed to execute subject template: " + err.Error())
		return subjectTemplate
	}

	result := buf.String()

	// Trim to max length if needed
	if len(result) > SubjectMaxLength {
		result = result[:SubjectMaxLength-3] + "..."
	}

	return strings.TrimSpace(result)
}

// GetTemplateByID finds a template by its ID
func GetTemplateByID(templateID string) (*MessageTemplate, error) {
	// Check connection templates
	for _, template := range GetConnectionRequestTemplates() {
		if template.ID == templateID {
			return &template, nil
		}
	}

	// Check message templates
	for _, template := range GetMessageTemplates() {
		if template.ID == templateID {
			return &template, nil
		}
	}

	return nil, fmt.Errorf("template not found: %s", templateID)
}

// GetTemplatesByType returns all templates of a specific type
func GetTemplatesByType(templateType TemplateType) []MessageTemplate {
	var templates []MessageTemplate

	if templateType == TemplateConnectionRequest {
		templates = GetConnectionRequestTemplates()
	} else {
		for _, template := range GetMessageTemplates() {
			if template.Type == templateType {
				templates = append(templates, template)
			}
		}
	}

	return templates
}

// ValidateMessageLength checks if a message is within LinkedIn's limits
func ValidateMessageLength(message string, messageType TemplateType) error {
	length := len(message)

	if messageType == TemplateConnectionRequest {
		if length > ConnectionNoteMaxLength {
			return fmt.Errorf("connection note too long: %d characters (max %d)", length, ConnectionNoteMaxLength)
		}
	} else {
		if length > MessageMaxLength {
			return fmt.Errorf("message too long: %d characters (max %d)", length, MessageMaxLength)
		}
	}

	if length == 0 {
		return fmt.Errorf("message cannot be empty")
	}

	return nil
}

// cleanupWhitespace removes excessive whitespace from text
func cleanupWhitespace(text string) string {
	// Replace multiple spaces with single space
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	// Replace multiple newlines with double newline
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}

	// Trim leading/trailing whitespace from each line
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	text = strings.Join(lines, "\n")

	return strings.TrimSpace(text)
}

// TruncateMessage truncates a message to fit within the specified length
func TruncateMessage(message string, maxLength int) string {
	if len(message) <= maxLength {
		return message
	}

	// Truncate with ellipsis
	return message[:maxLength-3] + "..."
}

package automation

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-rod/rod"

	"linkedin-automation/internal/logger"
	"linkedin-automation/internal/storage"
)

// ProcessDailyFollowUps handles the daily follow-up messaging workflow
func ProcessDailyFollowUps(page *rod.Page, db *storage.Database, rateLimiter *RateLimiter) error {
	logger.Info("Starting daily follow-up workflow...")

	// 1. Check for new connections (mark as accepted)
	if os.Getenv("CHECK_CONNECTION_STATUS") == "true" {
		if err := CheckRecentConnections(page, db); err != nil {
			logger.Error("Failed to check recent connections: " + err.Error())
		}
	}

	// 2. Check for replies (stop automation for them)
	if err := CheckInboxForReplies(page, db); err != nil {
		logger.Error("Failed to check inbox for replies: " + err.Error())
	}

	// 3. Send follow-up messages
	if os.Getenv("ENABLE_MESSAGING") == "true" {
		// Check rate limit
		if err := rateLimiter.CheckDailyLimit(TaskMessage); err != nil {
			logger.Warning("Messaging rate limit reached - skipping messages")
			return nil
		}

		maxMessages := 3
		if os.Getenv("MAX_MESSAGES_PER_RUN") != "" {
			fmt.Sscanf(os.Getenv("MAX_MESSAGES_PER_RUN"), "%d", &maxMessages)
		}

		profiles, err := db.GetAcceptedConnectionProfiles(maxMessages, 30)
		if err != nil {
			return fmt.Errorf("failed to get profiles for messaging: %w", err)
		}

		logger.Info(fmt.Sprintf("Found %d profiles for potential follow-up", len(profiles)))

		templateID := os.Getenv("MESSAGE_TEMPLATE")
		if templateID == "" {
			templateID = "msg_introduction"
		}

		for _, profile := range profiles {
			// Check rate limit again
			if err := rateLimiter.CheckDailyLimit(TaskMessage); err != nil {
				break
			}

			tmpl, err := GetTemplateByID(templateID)
			if err != nil {
				logger.Error("Template not found: " + err.Error())
				continue
			}

			firstName := profile.Name
			if parts := strings.Split(profile.Name, " "); len(parts) > 0 {
				firstName = parts[0]
			}

			vars := TemplateVariables{
				FirstName:    firstName,
				FullName:     profile.Name,
				Company:      profile.Company,
				Title:        profile.Title,
				YourName:     os.Getenv("YOUR_NAME"),
				YourTitle:    os.Getenv("YOUR_TITLE"),
				YourCompany:  os.Getenv("YOUR_COMPANY"),
				Industry:     os.Getenv("YOUR_INDUSTRY"),
				CustomReason: os.Getenv("MESSAGE_CUSTOM_REASON"),
			}

			body, err := RenderTemplate(*tmpl, vars)
			if err != nil {
				logger.Error("Failed to render template: " + err.Error())
				continue
			}

			req := MessageRequest{
				ProfileID:  profile.ID,
				ProfileURL: profile.ProfileURL,
				Name:       profile.Name,
				Body:       body,
				TemplateID: tmpl.ID,
			}

			if err := SendMessage(page, db, req); err != nil {
				logger.Error(fmt.Sprintf("Failed to send message to %s: %s", profile.Name, err.Error()))
			} else {
				rateLimiter.RecordAction(TaskMessage)
			}
		}
	}

	return nil
}

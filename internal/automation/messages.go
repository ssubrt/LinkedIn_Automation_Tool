package automation

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"

	"linkedin-automation/internal/logger"
	"linkedin-automation/internal/stealth"
	"linkedin-automation/internal/storage"
)

// SendMessage sends a direct message to a connection
func SendMessage(page *rod.Page, db *storage.Database, request MessageRequest) error {
	logger.Info(fmt.Sprintf("Sending message to: %s (%s)", request.Name, request.ProfileID))

	// Navigate to profile page
	logger.Info("Navigating to profile: " + request.ProfileURL)
	err := page.Navigate(request.ProfileURL)
	if err != nil {
		return fmt.Errorf("failed to navigate to profile: %w", err)
	}

	page.MustWaitLoad()
	stealth.RandomDelay(2000, 3000)

	// Click Message button
	logger.Info("Looking for Message button...")
	// Selectors for Message button
	messageSelectors := []string{
		"button[aria-label^='Message']",
		".pvs-profile-actions__action button:has-text('Message')",
		"button.artdeco-button--primary:has-text('Message')",
		"a[href^='/messaging/thread']", // Sometimes it's a link
	}

	var messageButton *rod.Element
	for _, sel := range messageSelectors {
		btn, err := page.Timeout(2 * time.Second).Element(sel)
		if err == nil && btn != nil {
			if visible, _ := btn.Visible(); visible {
				messageButton = btn
				break
			}
		}
	}

	if messageButton == nil {
		return fmt.Errorf("message button not found")
	}

	messageButton.Click(proto.InputMouseButtonLeft, 1)
	stealth.RandomDelay(1500, 2500)

	// Wait for message box to open
	// It might be a popup or a separate page. Usually a popup on the bottom right or overlay.
	// We look for the message input area.
	inputSelector := "div[role='textbox'][aria-label^='Write a message']"
	input, err := page.Timeout(5 * time.Second).Element(inputSelector)
	if err != nil {
		// Try alternative selector
		input, err = page.Timeout(2 * time.Second).Element(".msg-form__contenteditable")
		if err != nil {
			return fmt.Errorf("message input field not found: %w", err)
		}
	}

	// Type Body
	logger.Info("Typing message...")
	input.Input(request.Body)
	stealth.RandomDelay(1000, 2000)

	// Click Send
	sendButtonSelector := "button[type='submit']"
	sendButton, err := page.Timeout(3 * time.Second).Element(sendButtonSelector)
	if err != nil {
		// Try finding by text
		sendButton, err = page.Timeout(3*time.Second).ElementR("button", `\bSend\b`)
		if err != nil {
			return fmt.Errorf("send button not found")
		}
	}

	// Ensure button is clickable
	if visible, _ := sendButton.Visible(); !visible {
		return fmt.Errorf("send button not visible")
	}

	sendButton.Click(proto.InputMouseButtonLeft, 1)
	logger.Info("Message sent successfully")

	// Record in DB
	msg := storage.Message{
		ConnectionID:   request.ProfileID,
		TemplateName:   request.TemplateID,
		MessageContent: request.Body,
		SentAt:         time.Now(),
		CreatedAt:      time.Now(),
	}

	if err := db.SaveMessage(msg); err != nil {
		logger.Error("Failed to save message to database: " + err.Error())
	}

	return nil
}

package automation

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"

	"linkedin-automation/internal/logger"
	"linkedin-automation/internal/stealth"
	"linkedin-automation/internal/storage"
	"linkedin-automation/pkg/utils"
)

// CheckInboxForReplies checks the inbox for new replies and updates the database
func CheckInboxForReplies(page *rod.Page, db *storage.Database) error {
	logger.Info("Checking inbox for replies...")

	// Navigate to messaging
	err := page.Navigate("https://www.linkedin.com/messaging/")
	if err != nil {
		return fmt.Errorf("failed to navigate to messaging: %w", err)
	}

	page.MustWaitLoad()
	stealth.RandomDelay(2000, 3000)

	// Get list of conversations
	conversationSelector := ".msg-conversation-listitem"
	conversations, err := page.Timeout(5 * time.Second).Elements(conversationSelector)
	if err != nil {
		logger.Warning("Failed to get conversations or inbox empty: " + err.Error())
		return nil
	}

	logger.Info(fmt.Sprintf("Found %d conversations", len(conversations)))

	// Limit to top 10 to avoid excessive processing
	limit := 10
	if len(conversations) < limit {
		limit = len(conversations)
	}

	for i := 0; i < limit; i++ {
		// Re-fetch conversations to avoid stale elements
		conversations, _ = page.Elements(conversationSelector)
		if i >= len(conversations) {
			break
		}
		conv := conversations[i]

		// Click to open conversation
		conv.Click(proto.InputMouseButtonLeft, 1)
		stealth.RandomDelay(1000, 1500)

		// Identify the other person
		headerLink, err := page.Timeout(3 * time.Second).Element(".msg-entity-lockup__link")
		if err != nil {
			continue
		}

		href, err := headerLink.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		profileID := utils.ExtractProfileID(*href)
		if profileID == "" {
			continue
		}

		// Check last message
		bubbles, err := page.Elements(".msg-s-message-list__event")
		if err != nil || len(bubbles) == 0 {
			continue
		}

		lastBubble := bubbles[len(bubbles)-1]

		// Check if it is from me
		// LinkedIn uses classes like 'msg-s-message-list__event--s-me' for sent messages
		// and 'msg-s-message-list__event--other' for received messages.
		class, err := lastBubble.Attribute("class")
		if err != nil || class == nil {
			continue
		}

		isFromMe := strings.Contains(*class, "--s-me")

		if !isFromMe {
			// It's a reply!
			logger.Info(fmt.Sprintf("Detected reply from %s", profileID))
			err = db.UpdateConnectionReplyStatus(profileID, true)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to update reply status for %s: %s", profileID, err.Error()))
			}
		}
	}

	return nil
}

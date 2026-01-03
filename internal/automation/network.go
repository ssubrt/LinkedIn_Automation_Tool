package automation

import (
	"fmt"

	"github.com/go-rod/rod"

	"linkedin-automation/internal/logger"
	"linkedin-automation/internal/stealth"
	"linkedin-automation/internal/storage"
	"linkedin-automation/pkg/utils"
)

// CheckRecentConnections scrapes the "Recently Added" connections and updates the database
func CheckRecentConnections(page *rod.Page, db *storage.Database) error {
	logger.Info("Checking recent connections...")

	// Navigate to connections page
	err := page.Navigate("https://www.linkedin.com/mynetwork/invite-connect/connections/")
	if err != nil {
		return fmt.Errorf("failed to navigate to connections: %w", err)
	}

	page.MustWaitLoad()
	stealth.RandomDelay(2000, 3000)

	// Scrape the list
	// Selector for connection cards
	cardSelector := ".mn-connection-card"
	cards, err := page.Elements(cardSelector)
	if err != nil {
		return fmt.Errorf("failed to get connection cards: %w", err)
	}

	logger.Info(fmt.Sprintf("Found %d recent connections", len(cards)))

	count := 0
	for _, card := range cards {
		// Extract Profile URL
		link, err := card.Element("a.mn-connection-card__link")
		if err != nil {
			continue
		}

		href, err := link.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		profileID := utils.ExtractProfileID(*href)
		if profileID == "" {
			continue
		}

		// Update status to accepted
		err = db.UpdateConnectionStatus(profileID, "accepted")
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to update status for %s: %s", profileID, err.Error()))
		} else {
			count++
		}
	}

	logger.Info(fmt.Sprintf("Processed %d connections", count))
	return nil
}

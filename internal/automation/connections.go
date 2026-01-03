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

// ConnectionRequest represents a connection request to be sent
type ConnectionRequest struct {
	ProfileID   string
	ProfileURL  string
	Name        string
	Title       string
	Company     string
	Note        string
	TemplateID  string
	RequestedAt time.Time
}

// MessageRequest represents a message to be sent
type MessageRequest struct {
	ProfileID  string
	ProfileURL string
	Name       string
	Subject    string
	Body       string
	TemplateID string
	SentAt     time.Time
}

// ConnectionStats tracks statistics for connection requests
type ConnectionStats struct {
	TotalAttempted   int
	Successful       int
	Failed           int
	AlreadyConnected int
	Pending          int // Track pending connections separately
	Errors           []string
	StartTime        time.Time
	EndTime          time.Time
}

// MessagingStats tracks statistics for messages sent
type MessagingStats struct {
	TotalAttempted int
	Successful     int
	Failed         int
	Errors         []string
	StartTime      time.Time
	EndTime        time.Time
}

// SendConnectionRequest sends a connection request to a LinkedIn profile
//
// Edge Cases Handled:
// 1. Already Connected - Checks for "Connected" status and returns specific error
// 2. Already Pending - Checks for "Pending" status and returns specific error (not counted as failure)
// 3. 3rd-Degree Connections - If Connect button not visible, clicks "More..." dropdown to find it
// 4. Note Addition - Adds personalized note if provided and textarea is available
//
// Returns:
// - nil if connection request sent successfully
// - error with "already connected" if already connected
// - error with "connection pending" if request already pending
// - error if Connect button not found even in More... dropdown
func SendConnectionRequest(page *rod.Page, db *storage.Database, request ConnectionRequest) error {
	logger.Info(fmt.Sprintf("Sending connection request to: %s (%s)", request.Name, request.ProfileID))

	// Navigate to profile page
	logger.Info("Navigating to profile: " + request.ProfileURL)
	err := page.Navigate(request.ProfileURL)
	if err != nil {
		return fmt.Errorf("failed to navigate to profile: %w", err)
	}

	page.MustWaitLoad()

	// Check for LinkedIn checkpoint/verification page
	currentURL := page.MustInfo().URL
	if utils.IsLinkedInCheckpoint(currentURL) {
		logger.Error("❌ LinkedIn checkpoint/verification detected at: " + currentURL)
		return fmt.Errorf("linkedin checkpoint detected, manual verification required")
	}
	stealth.RandomDelay(2000, 3000)

	// Apply random scroll to simulate reading profile
	stealth.RandomScroll(page)
	stealth.RandomDelay(1000, 2000)

	// Check if already connected
	// Use Timeout to avoid hanging if element doesn't exist
	alreadyConnectedMessage, _ := page.Timeout(2 * time.Second).Element(utils.AlreadyConnectedSelector)
	if alreadyConnectedMessage != nil {
		logger.Info("Already connected with " + request.Name)
		return fmt.Errorf("already connected")
	}

	// Check if connection request is pending
	pendingMessage, _ := page.Timeout(2 * time.Second).Element(utils.PendingConnectionSelector)
	if pendingMessage != nil {
		logger.Info("Connection request already pending for " + request.Name)
		return fmt.Errorf("connection pending")
	}

	// Look for "Connect" button
	// IMPORTANT: We must avoid sidebar suggestions and only act on
	// the primary profile header. To do this we scope our searches
	// to the <main> content area and, when possible, the
	// `.pvs-profile-actions` toolbar.
	var connectButton *rod.Element
	var found bool

	// Find main content container
	var mainEl *rod.Element
	mainEl, _ = page.Timeout(3 * time.Second).Element("main")

	// Strategy 1: Look inside the profile actions toolbar
	if mainEl != nil {
		logger.Info("Strategy 1: Searching for Connect button in main profile actions bar...")
		actionsEl, _ := mainEl.Element(".pvs-profile-actions")
		if actionsEl != nil {
			// Try text-based search first
			btn, err := actionsEl.ElementR("button", `\bConnect\b`)
			if err == nil && btn != nil {
				if visible, _ := btn.Visible(); visible {
					connectButton = btn
					found = true
				}
			}

			// Fallback to selector-based search inside actions bar
			if !found {
				selectors := []string{
					utils.ConnectButtonSelector,
					utils.ConnectButtonAltSelector,
					"button[aria-label='Connect']",
					"button[aria-label='Invite to connect']",
				}

				for _, sel := range selectors {
					btn, err := actionsEl.Element(sel)
					if err == nil && btn != nil {
						if visible, _ := btn.Visible(); visible {
							logger.Info("Found Connect button by selector in actions bar: " + sel)
							connectButton = btn
							found = true
							break
						}
					}
				}
			}
		}
	}

	// Strategy 2: Fallback to searching within <main> only (still avoids sidebar)
	if !found && mainEl != nil {
		logger.Info("Strategy 2: Searching for Connect button within <main>...")
		btn, err := mainEl.ElementR("button", `\bConnect\b`)
		if err == nil && btn != nil {
			if visible, _ := btn.Visible(); visible {
				logger.Info("Found Connect button by text within <main>")
				connectButton = btn
				found = true
			}
		}
	}

	// Strategy 3: Check "More" dropdown (scoped to main/profile header only)
	if !found {
		logger.Info("Connect button not found directly. Checking 'More' dropdown in main profile area...")

		var moreButton *rod.Element

		// Prefer searching for More inside the profile actions bar, then within <main>.
		var moreSearchRoots []*rod.Element
		if mainEl != nil {
			if actionsEl, _ := mainEl.Element(".pvs-profile-actions"); actionsEl != nil {
				moreSearchRoots = append(moreSearchRoots, actionsEl)
			}
			moreSearchRoots = append(moreSearchRoots, mainEl)
		}

		moreSelectors := []string{
			utils.MoreActionsButtonSelector,
			utils.MoreActionsButtonAltSelector,
			"button[aria-label='More actions']",
			"button:has-text('More')",
		}

		for _, root := range moreSearchRoots {
			for _, sel := range moreSelectors {
				btn, err := root.Timeout(1 * time.Second).Element(sel)
				if err == nil && btn != nil {
					text, _ := btn.Text()
					aria, _ := btn.Attribute("aria-label")
					if strings.Contains(text, "More") || (aria != nil && strings.Contains(*aria, "More")) {
						if visible, _ := btn.Visible(); visible {
							logger.Info("Found More button in main/profile header with selector: " + sel)
							moreButton = btn
							break
						}
					}
				}
			}
			if moreButton != nil {
				break
			}
		}

		// As a very last resort (should rarely be needed), allow a page-wide search
		if moreButton == nil {
			for _, sel := range moreSelectors {
				btn, err := page.Timeout(1 * time.Second).Element(sel)
				if err == nil && btn != nil {
					text, _ := btn.Text()
					aria, _ := btn.Attribute("aria-label")
					if strings.Contains(text, "More") || (aria != nil && strings.Contains(*aria, "More")) {
						if visible, _ := btn.Visible(); visible {
							logger.Info("Fallback: Found More button with page-wide search and selector: " + sel)
							moreButton = btn
							break
						}
					}
				}
			}
		}

		if moreButton != nil {
			logger.Info("Clicking More... button")
			moreButton.ScrollIntoView()
			stealth.RandomDelay(500, 1000)
			moreButton.Click(proto.InputMouseButtonLeft, 1)
			stealth.RandomDelay(800, 1200)

			// After clicking More, a dropdown menu with role="menu" should appear
			logger.Info("Waiting for dropdown menu after clicking More...")
			menuEl, err := page.Timeout(3 * time.Second).Element("div[role='menu']")
			if err != nil || menuEl == nil {
				logger.Warning("Dropdown menu not found after clicking More")
			} else {
				// Look for any element inside the menu whose visible text contains 'Connect'
				logger.Info("Searching for 'Connect' item inside dropdown menu...")
				btn, err := menuEl.ElementR("*", `(?i)\\bConnect\\b`)
				if err == nil && btn != nil {
					if visible, _ := btn.Visible(); visible {
						logger.Info("Found Connect item inside dropdown menu")
						connectButton = btn
						found = true
					}
				} else {
					logger.Warning("No 'Connect' item found inside dropdown menu")
				}
			}
		}
	}

	// If still no Connect button found
	if !found || connectButton == nil {
		// Before giving up, check if this profile is likely
		// already connected: presence of a primary Message button
		// without any Connect option.
		logger.Info("Connect button not found, checking if profile is already connected...")
		msgButton, _ := page.Timeout(2 * time.Second).Element(utils.MessageButtonSelector)
		if msgButton == nil {
			msgButton, _ = page.Timeout(2 * time.Second).Element(utils.MessageButtonAltSelector)
		}
		if msgButton != nil {
			if visible, _ := msgButton.Visible(); visible {
				logger.Info("Message button present but no Connect button - treating as already connected")
				return fmt.Errorf("already connected")
			}
		}

		return fmt.Errorf("connect button not found - profile may be out of network")
	}

	// Scroll button into view
	err = connectButton.ScrollIntoView()
	if err != nil {
		return fmt.Errorf("failed to scroll connect button into view: %w", err)
	}

	stealth.RandomDelay(500, 1000)

	// Click Connect button
	logger.Info("Clicking Connect button...")
	err = connectButton.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("failed to click connect button: %w", err)
	}

	stealth.RandomDelay(1500, 2500)
	// Wait for modal to appear (don't use MustWaitLoad as it might not trigger a full page load)

	// Check if "Add a note" modal appeared
	// We need to wait a bit for the modal animation
	time.Sleep(2 * time.Second)

	// Check for modal presence
	_, err = page.Timeout(5 * time.Second).Element(".artdeco-modal")
	if err != nil {
		logger.Warning("Modal did not appear after clicking Connect. Checking if request was sent automatically...")
	}

	if request.Note != "" {
		logger.Info("Adding personalized note...")

		// Look for "Add a note" button
		addNoteButton, _ := page.Timeout(3 * time.Second).Element(utils.AddNoteButtonSelector)
		if addNoteButton == nil {
			// Try finding by text
			addNoteButton, _ = page.Timeout(3*time.Second).ElementR("button", "Add a note")
		}

		if addNoteButton != nil {
			// Click "Add a note" button
			err = addNoteButton.Click(proto.InputMouseButtonLeft, 1)
			if err != nil {
				logger.Warning("Failed to click Add Note button: " + err.Error())
			} else {
				stealth.RandomDelay(1000, 1500)

				// Find the note textarea
				noteTextarea, err := page.Timeout(3 * time.Second).Element(utils.ConnectionNoteTextareaSelector)
				if err != nil || noteTextarea == nil {
					noteTextarea, err = page.Timeout(3 * time.Second).Element("textarea[name='message']")
				}

				if err == nil && noteTextarea != nil {
					// Remove timeout context from the element for long operations like typing
					noteTextarea = noteTextarea.CancelTimeout()

					// Type the note with human-like typing
					logger.Info(fmt.Sprintf("Typing note (%d characters)...", len(request.Note)))
					stealth.TypeLikeHuman(noteTextarea, request.Note)
					stealth.RandomDelay(1000, 2000)
				} else {
					logger.Warning("Note textarea not found")
				}
			}
		} else {
			logger.Warning("Add a note button not found, skipping note.")
		}
	}

	// Find and click the "Send" button
	logger.Info("Looking for Send button...")
	var sendButton *rod.Element

	// Selectors for Send button
	sendSelectors := []string{
		utils.SendConnectionButtonSelector,
		"button[aria-label='Send now']",
		"button[aria-label='Send invitation']",
		"button.artdeco-button--primary:has-text('Send')",
		"button:has-text('Send without a note')", // Fallback if note failed
	}

	for _, sel := range sendSelectors {
		btn, err := page.Timeout(2 * time.Second).Element(sel)
		if err == nil && btn != nil {
			if visible, _ := btn.Visible(); visible {
				sendButton = btn
				break
			}
		}
	}

	if sendButton == nil {
		// Try finding by text regex as last resort
		sendButton, _ = page.Timeout(2*time.Second).ElementR("button", `\bSend\b`)
	}

	if sendButton == nil {
		return fmt.Errorf("send button not found")
	}

	stealth.RandomDelay(500, 1000)

	logger.Info("Clicking Send button...")
	err = sendButton.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("failed to click send button: %w", err)
	}

	stealth.RandomDelay(2000, 3000)
	page.MustWaitLoad()

	// Save to database
	if db != nil {
		connectionReq := storage.ConnectionRequest{
			ProfileID: request.ProfileID,
			SentAt:    time.Now(),
			NoteUsed:  request.Note,
			Status:    "pending",
		}

		err = db.SaveConnectionRequest(connectionReq)
		if err != nil {
			logger.Warning("Failed to save connection request to database: " + err.Error())
		}
	}

	logger.Info("Connection request sent successfully to " + request.Name)
	return nil
}

// SendConnectionRequests sends multiple connection requests with rate limiting
func SendConnectionRequests(page *rod.Page, db *storage.Database, rateLimiter *RateLimiter, requests []ConnectionRequest) *ConnectionStats {
	stats := &ConnectionStats{
		StartTime: time.Now(),
	}

	logger.Info(fmt.Sprintf("Sending %d connection requests...", len(requests)))

	for _, request := range requests {
		stats.TotalAttempted++

		// Check rate limit
		err := rateLimiter.CheckDailyLimit(TaskConnection)
		if err != nil {
			logger.Warning("Connection rate limit reached: " + err.Error())
			stats.Errors = append(stats.Errors, "Rate limit reached")
			break
		}

		// Send the request
		err = SendConnectionRequest(page, db, request)
		if err != nil {
			if strings.Contains(err.Error(), "already connected") {
				stats.AlreadyConnected++
			} else if strings.Contains(err.Error(), "connection pending") {
				stats.Pending++
				logger.Info(fmt.Sprintf("Connection request already pending for %s", request.Name))
			} else {
				stats.Failed++
				stats.Errors = append(stats.Errors, fmt.Sprintf("%s: %s", request.Name, err.Error()))
				logger.Warning(fmt.Sprintf("Failed to send connection to %s: %s", request.Name, err.Error()))
			}
		} else {
			stats.Successful++

			// Record action for rate limiting
			if err := rateLimiter.RecordAction(TaskConnection); err != nil {
				logger.Warning("Failed to record connection action: " + err.Error())
			}
		}

		// Apply cooldown between connections
		if stats.TotalAttempted < len(requests) {
			rateLimiter.ApplyCooldown()
		}
	}

	stats.EndTime = time.Now()
	duration := stats.EndTime.Sub(stats.StartTime)

	logger.Info(fmt.Sprintf("Connection requests completed: %d successful, %d failed, %d already connected in %s",
		stats.Successful, stats.Failed, stats.AlreadyConnected, duration))

	return stats
}

// SendMessage function has been moved to messages.go

// SendMessages sends multiple messages with rate limiting
func SendMessages(page *rod.Page, db *storage.Database, rateLimiter *RateLimiter, messages []MessageRequest) *MessagingStats {
	stats := &MessagingStats{
		StartTime: time.Now(),
	}

	logger.Info(fmt.Sprintf("Sending %d messages...", len(messages)))

	for _, message := range messages {
		stats.TotalAttempted++

		// Check rate limit
		err := rateLimiter.CheckDailyLimit(TaskMessage)
		if err != nil {
			logger.Warning("Messaging rate limit reached: " + err.Error())
			stats.Errors = append(stats.Errors, "Rate limit reached")
			break
		}

		// Send the message
		err = SendMessage(page, db, message)
		if err != nil {
			stats.Failed++
			stats.Errors = append(stats.Errors, fmt.Sprintf("%s: %s", message.Name, err.Error()))
			logger.Warning(fmt.Sprintf("Failed to send message to %s: %s", message.Name, err.Error()))
		} else {
			stats.Successful++

			// Record action for rate limiting
			if err := rateLimiter.RecordAction(TaskMessage); err != nil {
				logger.Warning("Failed to record message action: " + err.Error())
			}
		}

		// Apply cooldown between messages
		if stats.TotalAttempted < len(messages) {
			rateLimiter.ApplyCooldown()
		}
	}

	stats.EndTime = time.Now()
	duration := stats.EndTime.Sub(stats.StartTime)

	logger.Info(fmt.Sprintf("Messaging completed: %d successful, %d failed in %s",
		stats.Successful, stats.Failed, duration))

	return stats
}

// PrepareConnectionRequestFromProfile creates a ConnectionRequest from a database profile
func PrepareConnectionRequestFromProfile(profile storage.Profile, templateID string, senderVars TemplateVariables) (*ConnectionRequest, error) {
	// Get template
	template, err := GetTemplateByID(templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	if template.Type != TemplateConnectionRequest {
		return nil, fmt.Errorf("template %s is not a connection request template", templateID)
	}

	// Prepare template variables
	vars := TemplateVariables{
		FullName:     profile.Name,
		Title:        profile.Title,
		Company:      profile.Company,
		YourName:     senderVars.YourName,
		YourTitle:    senderVars.YourTitle,
		YourCompany:  senderVars.YourCompany,
		CustomReason: senderVars.CustomReason,
		Industry:     senderVars.Industry,
	}

	// Extract first name
	parts := strings.Split(profile.Name, " ")
	if len(parts) > 0 {
		vars.FirstName = parts[0]
		if len(parts) > 1 {
			vars.LastName = strings.Join(parts[1:], " ")
		}
	}

	// Render the template
	note, err := RenderTemplate(*template, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	// Validate length
	if err := ValidateMessageLength(note, TemplateConnectionRequest); err != nil {
		return nil, err
	}

	return &ConnectionRequest{
		ProfileID:   profile.ID,
		ProfileURL:  profile.ProfileURL,
		Name:        profile.Name,
		Title:       profile.Title,
		Company:     profile.Company,
		Note:        note,
		TemplateID:  templateID,
		RequestedAt: time.Now(),
	}, nil
}

// PrepareMessageFromProfile creates a MessageRequest from a database profile
func PrepareMessageFromProfile(profile storage.Profile, templateID string, senderVars TemplateVariables) (*MessageRequest, error) {
	// Get template
	template, err := GetTemplateByID(templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	if template.Type == TemplateConnectionRequest {
		return nil, fmt.Errorf("template %s is a connection request template, not a message template", templateID)
	}

	// Prepare template variables
	vars := TemplateVariables{
		FullName:     profile.Name,
		Title:        profile.Title,
		Company:      profile.Company,
		YourName:     senderVars.YourName,
		YourTitle:    senderVars.YourTitle,
		YourCompany:  senderVars.YourCompany,
		CustomReason: senderVars.CustomReason,
		Industry:     senderVars.Industry,
	}

	// Extract first name
	parts := strings.Split(profile.Name, " ")
	if len(parts) > 0 {
		vars.FirstName = parts[0]
		if len(parts) > 1 {
			vars.LastName = strings.Join(parts[1:], " ")
		}
	}

	// Render the template body
	body, err := RenderTemplate(*template, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	// Render the subject
	subject := RenderSubject(template.Subject, vars)

	// Validate length
	if err := ValidateMessageLength(body, template.Type); err != nil {
		return nil, err
	}

	return &MessageRequest{
		ProfileID:  profile.ID,
		ProfileURL: profile.ProfileURL,
		Name:       profile.Name,
		Subject:    subject,
		Body:       body,
		TemplateID: templateID,
		SentAt:     time.Now(),
	}, nil
}

// CheckAndUpdateConnectionStatuses checks pending connection requests and updates their status
// This function navigates to the "My Network" page to check which connections were accepted
func CheckAndUpdateConnectionStatuses(page *rod.Page, db *storage.Database) (int, error) {
	logger.Info("Checking connection request statuses...")

	// Navigate to My Network page
	err := page.Navigate("https://www.linkedin.com/mynetwork/")
	if err != nil {
		return 0, fmt.Errorf("failed to navigate to My Network: %w", err)
	}

	page.MustWaitLoad()

	// Check for LinkedIn checkpoint/verification page
	currentURL := page.MustInfo().URL
	if utils.IsLinkedInCheckpoint(currentURL) {
		logger.Error("❌ LinkedIn checkpoint/verification detected at: " + currentURL)
		return 0, fmt.Errorf("linkedin checkpoint detected, manual verification required")
	}

	stealth.RandomDelay(2000, 3000)

	// Scroll to load content
	stealth.RandomScroll(page)
	stealth.RandomDelay(1000, 2000)

	// Get all pending connection requests from database
	pendingRequests, err := db.GetPendingConnections()
	if err != nil {
		return 0, fmt.Errorf("failed to get pending connections: %w", err)
	}

	if len(pendingRequests) == 0 {
		logger.Info("No pending connections to check")
		return 0, nil
	}

	logger.Info(fmt.Sprintf("Checking status for %d pending connections", len(pendingRequests)))

	acceptedCount := 0

	// For each pending connection, check if they're now in "My Network"
	for _, request := range pendingRequests {
		profileID := request.ProfileID
		// Navigate to their profile
		profileURL := fmt.Sprintf("https://www.linkedin.com/in/%s/", profileID)
		err := page.Navigate(profileURL)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to navigate to profile %s: %s", profileID, err.Error()))
			continue
		}

		page.MustWaitLoad()
		stealth.RandomDelay(1500, 2500)

		// Check for "Connected" indicator
		connectedElement, _ := page.Element(utils.AlreadyConnectedSelector)
		if connectedElement != nil {
			// Connection was accepted!
			logger.Info(fmt.Sprintf("Connection accepted: %s", profileID))
			err = db.UpdateConnectionStatus(profileID, "accepted")
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed to update status for %s: %s", profileID, err.Error()))
			} else {
				acceptedCount++
			}
		}

		// Apply delay between checks
		stealth.RandomDelay(2000, 3000)
	}

	logger.Info(fmt.Sprintf("Found %d newly accepted connections", acceptedCount))
	return acceptedCount, nil
}

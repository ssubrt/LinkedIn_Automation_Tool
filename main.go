package main

import (
	"fmt"
	"os"
	"time"

	"linkedin-automation/internal/automation"
	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/logger"
	"linkedin-automation/internal/stealth"
	"linkedin-automation/internal/storage"

	"github.com/go-rod/rod"
	"github.com/joho/godotenv"
)

// main orchestrates the LinkedIn automation workflow:
// 1. Loads environment variables
// 2. Checks activity scheduling (business hours only)
// 3. Initializes database and rate limiter
// 4. Checks for existing session
// 5. Initializes a browser instance with persistent session
// 6. Applies comprehensive fingerprint masking
// 7. Performs login only if needed
// 8. Executes advanced stealth actions
func main() {
	// Log the start of the automation process
	logger.Info("Starting LinkedIn Automation with Advanced Stealth")

	// Step 1: Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logger.Warning("No .env file found, using default configuration")
	}

	// Step 2: Check if we're in active hours (business hours)
	// logger.Info("Checking activity schedule...")
	// if !automation.IsActiveHours() {
	// 	logger.Warning("Outside active hours - waiting for business hours...")
	// 	automation.WaitForActiveHours()
	// }
	// logger.Info("Within active hours - proceeding with automation")

	// Step 3: Initialize SQLite database
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/linkedin_automation.db"
	}
	logger.Info("Initializing database at: " + dbPath)

	db, err := storage.InitDB(dbPath)
	if err != nil {
		logger.Error("Failed to initialize database: " + err.Error())
		return
	}
	defer db.Close()
	logger.Info("Database initialized successfully")

	// Step 3.5: Initialize rate limiter
	rateLimiter := automation.NewRateLimiter(db)

	// Display current rate limit stats
	stats, err := rateLimiter.GetDailyStats()
	if err != nil {
		logger.Warning("Failed to get rate limit stats: " + err.Error())
	} else {
		logger.Info("Rate Limiter initialized")
		fmt.Println(stats)
	}

	// Step 4: Check for existing session
	logger.Info("Checking for existing session...")
	state, err := storage.LoadState()
	if err != nil {
		logger.Warning("Failed to load state: " + err.Error())
	}

	sessionValid := false
	if state != nil && storage.IsSessionValid(state) {
		logger.Info("Valid session found! Skipping login...")
		sessionValid = true
	} else {
		logger.Info("No valid session found, login will be required")
	}

	// Step 5: Start the browser instance with persistent session support
	br, err := browser.StartBrowser()
	if err != nil {
		logger.Error("Failed to start Browser: " + err.Error())
		return
	}
	// Ensure browser is properly closed when the function exits
	defer br.Close()

	// Step 5.5: Apply comprehensive fingerprint masking BEFORE any page loads
	logger.Info("Applying advanced fingerprint masking...")
	browser.ApplyFingerprintMasking(br)

	// Step 6: Open LinkedIn and perform login if needed
	var page *rod.Page

	if sessionValid {
		// Try to navigate to LinkedIn home page directly
		logger.Info("Attempting to access LinkedIn with existing session...")
		page, err = browser.OpenPage(br, "https://www.linkedin.com/feed/")
		if err != nil {
			logger.Error("Failed to open LinkedIn: " + err.Error())
			return
		}

		// Wait a moment for page to load
		page.MustWaitLoad()

		// Check if we're actually logged in by checking the current URL
		currentURL := page.MustInfo().URL
		if currentURL == "https://www.linkedin.com/feed/" ||
			len(currentURL) > 0 && currentURL[:28] == "https://www.linkedin.com/feed" {
			logger.Info("Successfully accessed LinkedIn with saved session!")
		} else {
			// Session expired, need to login
			logger.Warning("Session expired, proceeding with login...")
			sessionValid = false

			// Navigate to login page
			page, err = browser.OpenPage(br, "https://www.linkedin.com/login")
			if err != nil {
				logger.Error("Failed to open LinkedIn login page: " + err.Error())
				return
			}
		}
	}

	if !sessionValid {
		// Open the LinkedIn login page
		page, err = browser.OpenPage(br, "https://www.linkedin.com/login")
		if err != nil {
			logger.Error("Failed to open LinkedIn Page: " + err.Error())
			return
		}

		// Read LinkedIn credentials from environment variables
		email := os.Getenv("LINKEDIN_EMAIL")
		password := os.Getenv("LINKEDIN_PASSWORD")

		if email == "" || password == "" {
			logger.Error("LINKEDIN_EMAIL or LINKEDIN_PASSWORD not set in .env file")
			return
		}

		// Perform the login action with credentials
		err = automation.LoginLinkedln(page, email, password)
		if err != nil {
			logger.Error("Login Failed: " + err.Error())
			// Invalidate session on failed login
			storage.InvalidateSession()
			return
		}
		logger.Info("Login Successful")

		// Save successful login state
		err = storage.SaveState(true)
		if err != nil {
			logger.Warning("Failed to save state: " + err.Error())
		}
	}

	// Step 7: Execute comprehensive stealth actions
	logger.Info("Starting advanced human-like behavior simulation...")

	// 7.1: Random mouse movements with Bézier curves
	logger.Info("Executing Bézier curve mouse movements...")
	stealth.MoveMouseRandomly(page)

	// 7.2: Hover over random elements (links, buttons)
	logger.Info("Hovering over interactive elements...")
	if err := stealth.HoverRandomElements(page); err != nil {
		logger.Warning("Failed to hover elements: " + err.Error())
	}

	// 7.3: Random scrolling with natural patterns
	logger.Info("Executing natural scrolling patterns...")
	stealth.RandomScroll(page)

	// Step 8: Execute LinkedIn people search
	logger.Info("Starting LinkedIn people search...")

	// Check rate limit before searching
	err = rateLimiter.CheckDailyLimit(automation.TaskSearch)
	canSearch := (err == nil)

	if canSearch {
		// Configure search parameters from environment variables
		searchConfig := automation.SearchConfig{
			Keywords:       os.Getenv("SEARCH_KEYWORDS"),
			JobTitle:       os.Getenv("SEARCH_JOB_TITLE"),
			Company:        os.Getenv("SEARCH_COMPANY"),
			Location:       os.Getenv("SEARCH_LOCATION"),
			MaxPages:       3, // Limit to 3 pages for now
			SkipDuplicates: true,
			DuplicateDays:  30,
		}

		// Use default values if environment variables are not set
		if searchConfig.Keywords == "" {
			searchConfig.Keywords = "software engineer"
		}
		if searchConfig.Location == "" {
			searchConfig.Location = "San Francisco Bay Area"
		}

		logger.Info("Search configuration:")
		logger.Info(fmt.Sprintf("  Keywords: %s", searchConfig.Keywords))
		logger.Info(fmt.Sprintf("  Job Title: %s", searchConfig.JobTitle))
		logger.Info(fmt.Sprintf("  Company: %s", searchConfig.Company))
		logger.Info(fmt.Sprintf("  Location: %s", searchConfig.Location))

		// Execute the search
		searchResults, searchStats, err := automation.SearchPeople(page, db, searchConfig)
		if err != nil {
			logger.Error("Search failed: " + err.Error())
		} else {
			// Record search action in rate limiter
			if err := rateLimiter.RecordAction(automation.TaskSearch); err != nil {
				logger.Warning("Failed to record search action: " + err.Error())
			}

			// Display search statistics
			logger.Info("Search completed successfully!")
			fmt.Println("\n========== Search Statistics ==========")
			fmt.Printf("Total profiles found: %d\n", searchStats.TotalFound)
			fmt.Printf("New profiles saved: %d\n", searchStats.NewProfiles)
			fmt.Printf("Duplicates skipped: %d\n", searchStats.Duplicates)
			fmt.Printf("Pages scraped: %d\n", searchStats.PagesScraped)
			fmt.Printf("Errors encountered: %d\n", searchStats.ErrorCount)
			fmt.Printf("Duration: %s\n", searchStats.EndTime.Sub(searchStats.StartTime))
			fmt.Println("=======================================")

			// Warn if no profiles found - likely indicates selector changes
			if searchStats.TotalFound == 0 && searchStats.PagesScraped > 0 {
				logger.Warning("⚠️  Zero profiles found despite successful page load!")
				logger.Warning("⚠️  LinkedIn may have changed their HTML selectors.")
				logger.Warning("⚠️  Check constants.go and update SearchResultItemSelector if needed.")
			}

			// IMMEDIATE CONNECTION FLOW
			// Connect to found profiles immediately (limit to 3)
			if len(searchResults) > 0 && os.Getenv("ENABLE_CONNECTIONS") == "true" {
				logger.Info("Starting immediate connection requests for found profiles...")

				count := 0
				for _, result := range searchResults {
					if count >= 3 {
						break
					}

					// Check rate limit
					if err := rateLimiter.CheckDailyLimit(automation.TaskConnection); err != nil {
						logger.Warning("Connection rate limit reached")
						break
					}

					// Prepare request
					req := automation.ConnectionRequest{
						ProfileID:   result.ProfileID,
						ProfileURL:  result.ProfileURL,
						Name:        result.Name,
						Title:       result.Title,
						Company:     result.Company,
						RequestedAt: time.Now(),
						// TemplateID can be added here if needed
					}

					// Send request
					err := automation.SendConnectionRequest(page, db, req)
					if err != nil {
						logger.Error("Failed to connect to " + result.Name + ": " + err.Error())
					} else {
						logger.Info("Connection request sent to " + result.Name)
						rateLimiter.RecordAction(automation.TaskConnection)
						count++
					}
				}
			}
		}
	} else {
		logger.Warning("Search rate limit reached - skipping search for today")
	}

	// Step 9: Send connection requests (if enabled)
	// NOTE: This step is redundant if we are doing immediate connections above.
	// However, it's useful for processing profiles found in previous runs.
	if os.Getenv("ENABLE_CONNECTIONS") == "true" {
		logger.Info("Starting connection request automation (processing backlog)...")

		// Check rate limit
		err = rateLimiter.CheckDailyLimit(automation.TaskConnection)
		if err == nil {
			// Get profiles that haven't been contacted yet
			maxConnections := 5 // Limit to 5 connections per run for safety
			if os.Getenv("MAX_CONNECTIONS_PER_RUN") != "" {
				fmt.Sscanf(os.Getenv("MAX_CONNECTIONS_PER_RUN"), "%d", &maxConnections)
			}

			profiles, err := db.GetRecentProfiles(maxConnections, 30) // Get up to 5 profiles from last 30 days
			if err != nil {
				logger.Warning("Failed to get profiles for connections: " + err.Error())
			} else if len(profiles) > 0 {
				logger.Info(fmt.Sprintf("Found %d profiles for connection requests", len(profiles)))

				// Prepare sender variables from environment
				senderVars := automation.TemplateVariables{
					YourName:     os.Getenv("YOUR_NAME"),
					YourTitle:    os.Getenv("YOUR_TITLE"),
					YourCompany:  os.Getenv("YOUR_COMPANY"),
					Industry:     os.Getenv("YOUR_INDUSTRY"),
					CustomReason: os.Getenv("CONNECTION_CUSTOM_REASON"),
				}

				// Get template ID from environment (default to generic)
				templateID := os.Getenv("CONNECTION_TEMPLATE")
				if templateID == "" {
					templateID = "conn_generic"
				}

				// Prepare connection requests
				var requests []automation.ConnectionRequest
				for _, profile := range profiles {
					request, err := automation.PrepareConnectionRequestFromProfile(profile, templateID, senderVars)
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to prepare connection for %s: %s", profile.Name, err.Error()))
						continue
					}
					requests = append(requests, *request)
				}

				if len(requests) > 0 {
					// Send connection requests
					connStats := automation.SendConnectionRequests(page, db, rateLimiter, requests)

					// Display stats
					fmt.Println("\n========== Connection Request Statistics ==========")
					fmt.Printf("Total attempted: %d\n", connStats.TotalAttempted)
					fmt.Printf("Successful: %d\n", connStats.Successful)
					fmt.Printf("Failed: %d\n", connStats.Failed)
					fmt.Printf("Already connected: %d\n", connStats.AlreadyConnected)
					fmt.Printf("Already pending: %d\n", connStats.Pending)
					if len(connStats.Errors) > 0 {
						fmt.Printf("Errors: %d\n", len(connStats.Errors))
						for i, errMsg := range connStats.Errors {
							if i < 3 { // Show first 3 errors
								fmt.Printf("  - %s\n", errMsg)
							}
						}
					}
					fmt.Printf("Duration: %s\n", connStats.EndTime.Sub(connStats.StartTime))
					fmt.Println("===================================================")
				}
			} else {
				logger.Info("No profiles available for connection requests")
			}
		} else {
			logger.Warning("Connection rate limit reached - skipping connections for today")
		}
	}

	// Step 10: Execute daily follow-up workflow (Connection checks, Reply detection, Messaging)
	if os.Getenv("ENABLE_MESSAGING") == "true" || os.Getenv("CHECK_CONNECTION_STATUS") == "true" {
		err = automation.ProcessDailyFollowUps(page, db, rateLimiter)
		if err != nil {
			logger.Error("Daily follow-up workflow failed: " + err.Error())
		}
	}

	// Step 11: Display final stats
	logger.Info("Automation workflow completed successfully!")

	// Show rate limit summary
	if stats, err := rateLimiter.GetDailyStats(); err == nil {
		fmt.Println("\n" + stats)
	}

	logger.Info("Browser will remain open. Press Ctrl+C to exit.")

	// Keep the browser open to see results before closing
	select {}
}

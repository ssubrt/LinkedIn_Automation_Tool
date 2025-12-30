package main

import (
	"os"

	"linkedin-automation/internal/automation"
	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/logger"
	"linkedin-automation/internal/storage"

	"github.com/go-rod/rod"
	"github.com/joho/godotenv"
)

// main orchestrates the LinkedIn automation workflow:
// 1. Loads environment variables
// 2. Initializes database
// 3. Checks for existing session
// 4. Initializes a browser instance with persistent session
// 5. Performs login only if needed
// 6. Applies browser fingerprint masking for stealth
func main() {
	// Log the start of the automation process
	logger.Info("Starting LinkedIn Automation")

	// Step 1: Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logger.Warning("No .env file found, using default configuration")
	}

	// Step 2: Initialize SQLite database
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

	// Step 3: Check for existing session
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

	// Step 4: Start the browser instance with persistent session support
	br, err := browser.StartBrowser()
	if err != nil {
		logger.Error("Failed to start Browser: " + err.Error())
		return
	}
	// Ensure browser is properly closed when the function exits
	defer br.Close()

	// Step 5: Open LinkedIn and perform login if needed
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

	// Step 6: Apply browser fingerprint masking to avoid detection
	logger.Info("Applying fingerprint masking...")
	browser.ApplyFingerprintMasking(br)

	// Step 7: Perform stealth actions on the LinkedIn page
	logger.Info("Starting human-like behavior simulation...")
	browser.PerformStealthActions(page)

	logger.Info("Automation workflow completed successfully!")
	logger.Info("Browser will remain open. Press Ctrl+C to exit.")

	// Keep the browser open to see results before closing
	select {}
}

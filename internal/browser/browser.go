package browser

import (
	"fmt"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"linkedin-automation/internal/logger"
	"linkedin-automation/internal/stealth"
)

// BrowserConfig holds configuration for browser initialization
type BrowserConfig struct {
	UserDataDir string
	Headless    bool
}

// StartBrowser launches and returns a Rod Browser instance with persistent session support
func StartBrowser() (*rod.Browser, error) {
	return StartBrowserWithConfig(BrowserConfig{
		UserDataDir: "./browser_data",
		Headless:    false,
	})
}

// StartBrowserWithConfig launches a browser with custom configuration
func StartBrowserWithConfig(config BrowserConfig) (*rod.Browser, error) {
	logger.Info("Launching browser with persistent session...")

	// Ensure the user data directory exists
	if err := os.MkdirAll(config.UserDataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create user data directory: %w", err)
	}

	// Configure launcher with user data directory for session persistence
	l := launcher.New().
		Delete("leakless").
		Headless(config.Headless).
		UserDataDir(config.UserDataDir)

	u, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	logger.Info("Browser launched, connecting...")

	browser := rod.New().ControlURL(u)

	err = browser.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

	logger.Info("Browser connected successfully with persistent session!")

	return browser, nil
}

// PerformStealthActions executes human-like behavior on the page (mouse movements and scrolling)
// to avoid detection by anti-bot systems
func PerformStealthActions(page *rod.Page) {
	logger.Info("Performing stealth actions - simulating human-like behavior")

	// Perform random mouse movements
	logger.Info("Executing random mouse movements...")
	stealth.MoveMouseRandomly(page)

	// Perform random scrolling
	logger.Info("Executing random page scrolling...")
	stealth.RandomScroll(page)

	logger.Info("Stealth actions completed")
}

// OpenPage opens a new page, applies fingerprint masking, then navigates to the specified URL
// This prevents race conditions where detection scripts run before masking is applied
func OpenPage(browser *rod.Browser, url string) (*rod.Page, error) {
	// Create a blank page first
	page := browser.MustPage("about:blank")

	// CRITICAL: Apply fingerprint masking BEFORE navigation
	// This prevents LinkedIn's detection scripts from running before our masks are in place
	logger.Info("Applying fingerprint masking to page before navigation...")
	err := ApplyPageFingerprint(page)
	if err != nil {
		logger.Warning("Failed to apply fingerprint before navigation: " + err.Error())
		// Continue anyway - better to try with partial masking than fail completely
	}

	// NOW navigate to the target URL with masking already applied
	err = page.Navigate(url)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to %s: %w", url, err)
	}

	return page, nil
}

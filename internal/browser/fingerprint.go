package browser

import (
	"github.com/go-rod/rod"
)

// ApplyFingerprintMasking applies basic anti-detection measures to the browser.
// It configures the browser to mask itself as a legitimate user by ignoring certificate errors.
// User agent should be set during browser launch in the StartBrowser function.
func ApplyFingerprintMasking(br *rod.Browser) {
	// Ignore certificate transparency and HSTS (HTTP Strict Transport Security) errors
	// This prevents automation detection due to certificate validation issues
	// which is a common stealth technique used in web automation
	br.MustIgnoreCertErrors(true)
}

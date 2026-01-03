package browser

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"

	"linkedin-automation/internal/logger"
	"linkedin-automation/pkg/utils"
)

// ApplyFingerprintMasking applies comprehensive anti-detection measures to the browser.
func ApplyFingerprintMasking(br *rod.Browser) {
	// Ignore certificate errors
	br.MustIgnoreCertErrors(true)

	logger.Info("Applying advanced fingerprint masking...")

	// Get all pages and apply masking to each
	pages := br.MustPages()
	for _, page := range pages {
		if err := ApplyPageFingerprint(page); err != nil {
			logger.Warning("Failed to apply fingerprint to page: " + err.Error())
		}
	}
}

// ApplyPageFingerprint applies fingerprint masking to a specific page
func ApplyPageFingerprint(page *rod.Page) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// We construct a single large IIFE (Immediately Invoked Function Expression)
	// to ensure variables like 'const' don't leak or conflict, and comments don't break structure.

	// 1. Mask navigator.webdriver
	maskWebDriver := `
		try {
			Object.defineProperty(navigator, 'webdriver', {
				get: () => undefined
			});
		} catch (e) {}
	`

	// 2. Override automation-related properties
	maskAutomation := `
		try {
			// Remove automation indicators
			if (navigator.__proto__ && navigator.__proto__.webdriver) {
				delete navigator.__proto__.webdriver;
			}
			
			// Override chrome property
			if (!window.chrome) {
				window.chrome = {
					runtime: {},
					loadTimes: function() {},
					csi: function() {},
					app: {}
				};
			}
		} catch (e) {}
	`

	// 3. Randomize plugin array
	plugins := []string{
		"Chrome PDF Plugin",
		"Chrome PDF Viewer",
		"Native Client",
	}

	maskPlugins := fmt.Sprintf(`
		try {
			Object.defineProperty(navigator, 'plugins', {
				get: () => [
					{ name: '%s', filename: 'internal-pdf-viewer', description: 'Portable Document Format' },
					{ name: '%s', filename: 'mhjfbmdgcfjbbpaeojofohoefgiehjai', description: 'Portable Document Format' },
					{ name: '%s', filename: 'internal-nacl-plugin', description: 'Native Client Executable' }
				]
			});
		} catch (e) {}
	`, plugins[0], plugins[1], plugins[2])

	// 4. Randomize languages
	languages := []string{"en-US", "en"}
	maskLanguages := fmt.Sprintf(`
		try {
			Object.defineProperty(navigator, 'languages', {
				get: () => ['%s', '%s']
			});
		} catch (e) {}
	`, languages[0], languages[1])

	// 5. Override permissions API
	maskPermissions := `
		try {
			const originalQuery = window.navigator.permissions.query;
			window.navigator.permissions.query = (parameters) => (
				parameters.name === 'notifications' ?
					Promise.resolve({ state: Notification.permission }) :
					originalQuery(parameters)
			);
		} catch (e) {}
	`

	// 6. Mask canvas fingerprinting
	maskCanvas := `
		try {
			const originalGetImageData = CanvasRenderingContext2D.prototype.getImageData;
			
			const noisify = function(context, width, height, imageData) {
				const data = imageData.data;
				for (let i = 0; i < data.length; i += 4) {
					// Add minimal noise to RGB values
					const noise = Math.floor(Math.random() * 3) - 1;
					data[i] = data[i] + noise;
					data[i + 1] = data[i + 1] + noise;
					data[i + 2] = data[i + 2] + noise;
				}
				return imageData;
			};

		if (originalGetImageData) {
			Object.defineProperty(CanvasRenderingContext2D.prototype, 'getImageData', {
				value: function() {
					if (!originalGetImageData) return null;
					const imageData = originalGetImageData.apply(this, arguments);
					return noisify(this, arguments[2], arguments[3], imageData);
				}
			});
		}
	} catch (e) {}
`

	// 7. Mask WebGL fingerprinting
	maskWebGL := `
		try {
			const getParameter = WebGLRenderingContext.prototype.getParameter;
			WebGLRenderingContext.prototype.getParameter = function(parameter) {
				// Randomize vendor and renderer (UNMASKED_VENDOR_WEBGL = 37445, UNMASKED_RENDERER_WEBGL = 37446)
				if (parameter === 37445) {
					return 'Intel Inc.';
				}
				if (parameter === 37446) {
					return 'Intel Iris OpenGL Engine';
				}
				return getParameter.apply(this, arguments);
			};
		} catch (e) {}
	`

	// 8. Spoof screen properties
	screenWidth := 1920 + r.Intn(200) - 100  // 1820-2020
	screenHeight := 1080 + r.Intn(200) - 100 // 980-1180

	maskScreen := fmt.Sprintf(`
		try {
			Object.defineProperty(screen, 'width', { get: () => %d });
			Object.defineProperty(screen, 'height', { get: () => %d });
			Object.defineProperty(screen, 'availWidth', { get: () => %d });
			Object.defineProperty(screen, 'availHeight', { get: () => %d });
		} catch (e) {}
	`, screenWidth, screenHeight, screenWidth, screenHeight-40)

	// 9. Override battery API
	maskBattery := `
		try {
			if (navigator.getBattery) {
				navigator.getBattery = () => Promise.resolve({
					charging: true,
					chargingTime: 0,
					dischargingTime: Infinity,
					level: 1,
					addEventListener: () => {},
					removeEventListener: () => {},
					dispatchEvent: () => true
				});
			}
		} catch (e) {}
	`

	// 10. Mask connection API
	maskConnection := `
		try {
			if (navigator.connection) {
				Object.defineProperty(navigator, 'connection', {
					get: () => ({
						effectiveType: '4g',
						downlink: 10,
						rtt: 50,
						saveData: false
					})
				});
			}
		} catch (e) {}
	`

	// Combine all masking scripts inside an IIFE to isolate scope
	fullScript := fmt.Sprintf(`
		(function() {
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			%s
			%s
		})();
	`, maskWebDriver, maskAutomation, maskPlugins, maskLanguages,
		maskPermissions, maskCanvas, maskWebGL, maskScreen, maskBattery, maskConnection)

	// Execute the masking script
	_, err := page.Eval(fullScript)
	if err != nil {
		return fmt.Errorf("failed to apply fingerprint masking: %w", err)
	}

	// Set custom user agent
	err = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: utils.ChromeUserAgent,
	})
	if err != nil {
		return fmt.Errorf("failed to set user agent: %w", err)
	}

	// Randomize viewport size
	viewportWidth := 1366 + r.Intn(500) // 1366-1866
	viewportHeight := 768 + r.Intn(300) // 768-1068

	err = page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:             viewportWidth,
		Height:            viewportHeight,
		DeviceScaleFactor: 1,
		Mobile:            false,
	})
	if err != nil {
		return fmt.Errorf("failed to set viewport: %w", err)
	}

	logger.Info(fmt.Sprintf("Fingerprint applied: viewport %dx%d, screen %dx%d",
		viewportWidth, viewportHeight, screenWidth, screenHeight))

	return nil
}

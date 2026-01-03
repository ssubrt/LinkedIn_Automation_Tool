package automation

/*
	errors- used to return meaningful errors
	rod- browser interaction library
	stealth - humanlike behaviour functions
	logger - cleans logs instead of print

*/
import (
	"errors"
	"time"

	"github.com/go-rod/rod"

	"linkedin-automation/internal/logger"
	"linkedin-automation/internal/stealth"
)

/*
LoginLinkedln - logs into linkedin 	with given credentials
page - rod page to perform actions on (currently opened linkedin login page)
returns errors if any issue occurs during linkedin login
*/
func LoginLinkedln(page *rod.Page, email string, password string) error {

	//navigate to linkedin login page and wait until the page is fully loaded
	logger.Info("Opening Linkedin Login page")
	page.MustNavigate("https://www.linkedin.com/login")
	page.MustWaitLoad()

	//Human like delay between actions
	stealth.RandomDelay(1500, 3000)

	//Locate email iput field  and tell user if the email id input field is empty
	logger.Info("Locating email input field")
	emailInput, err := page.Timeout(10 * time.Second).Element(`input#username`)
	if err != nil {
		return errors.New("email input not found")
	}

	//pause for random time to mimic human behaviour and type email id like a human
	stealth.RandomDelay(800, 1500)
	emailInput.MustInput(email)

	//Locate password input field  and tell user if the password input field is empty
	logger.Info("Locating password input field")
	passwordInput, err := page.Timeout(10 * time.Second).Element(`input#password`)
	if err != nil {
		return errors.New("password input not found")
	}

	//pause for random time to mimic human behaviour and type password like a human
	stealth.RandomDelay(800, 1500)
	stealth.TypeLikeHuman(passwordInput, password)

	//Locate and click on the Sign in button to submit the credentials
	logger.Info("Locating and clicking on sign in button")
	loginBtn, err := page.Timeout(10 * time.Second).Element(`button[type="submit"]`)
	if err != nil {
		return errors.New("Login Button not found")
	}

	//puase for random time to show human behaviour and  click on the login button
	stealth.RandomDelay(1000, 2000)
	loginBtn.MustClick()

	//Wait for linkedln page to load after the login btn is clicked
	stealth.RandomDelay(3000, 5000)
	page.MustWaitLoad()

	// Check current URL first to see if login succeeded immediately
	logger.Info("Checking login status...")
	stealth.RandomDelay(2000, 3000)
	currentURL := page.MustInfo().URL
	logger.Info("Current page URL: " + currentURL)

	// If already on feed/home page, login succeeded without 2FA
	if currentURL != "https://www.linkedin.com/login" &&
		(len(currentURL) >= 28 && currentURL[:28] == "https://www.linkedin.com/feed" ||
			len(currentURL) >= 29 && currentURL[:29] == "https://www.linkedin.com/check") {
		logger.Info("✓ Login successful!")
		return nil
	}

	// Still on login page - check for challenges with timeout
	logger.Info("Checking for 2FA or CAPTCHA challenges...")

	// Check for 2FA challenge (with timeout)
	// twoFAChallenge, _ := page.Timeout(2 * time.Second).Element("#challenge")
	// if twoFAChallenge != nil {
	// 	logger.Warning("⚠️  2FA challenge detected! Manual intervention required.")
	// 	logger.Info("Please complete 2FA verification manually in the browser.")
	// 	logger.Info("Waiting up to 5 minutes for completion...")

	// 	// Wait up to 5 minutes, checking every 10 seconds
	// 	for i := 0; i < 30; i++ {
	// 		stealth.RandomDelay(10000, 10500)
	// 		currentURL = page.MustInfo().URL
	// 		if currentURL != "https://www.linkedin.com/login" {
	// 			logger.Info("✓ 2FA completed successfully!")
	// 			return nil
	// 		}
	// 		logger.Info(fmt.Sprintf("Still waiting... (%d/30 checks)", i+1))
	// 	}
	// 	return errors.New("2FA timeout - please try again")
	// }

	// Check for CAPTCHA (with timeout)
	// captchaChallenge, _ := page.Timeout(2 * time.Second).Element(".g-recaptcha")
	// if captchaChallenge != nil {
	// 	logger.Warning("⚠️  CAPTCHA challenge detected! Manual intervention required.")
	// 	logger.Info("Please complete CAPTCHA verification manually in the browser.")
	// 	logger.Info("Waiting up to 5 minutes for completion...")

	// 	// Wait up to 5 minutes, checking every 10 seconds
	// 	for i := 0; i < 30; i++ {
	// 		stealth.RandomDelay(10000, 10500)
	// 		currentURL = page.MustInfo().URL
	// 		if currentURL != "https://www.linkedin.com/login" {
	// 			logger.Info("✓ CAPTCHA completed successfully!")
	// 			return nil
	// 		}
	// 		logger.Info(fmt.Sprintf("Still waiting... (%d/30 checks)", i+1))
	// 	}
	// 	return errors.New("CAPTCHA timeout - please try again")
	// }

	// Check for security verification (with timeout)
	// securityChallenge, _ := page.Timeout(2 * time.Second).Element("form[action*='checkpoint']")
	// if securityChallenge != nil {
	// 	logger.Warning("⚠️  Security verification detected! Manual intervention required.")
	// 	logger.Info("Please complete security verification manually in the browser.")
	// 	logger.Info("Waiting up to 5 minutes for completion...")

	// 	// Wait up to 5 minutes, checking every 10 seconds
	// 	for i := 0; i < 30; i++ {
	// 		stealth.RandomDelay(10000, 10500)
	// 		currentURL = page.MustInfo().URL
	// 		if currentURL != "https://www.linkedin.com/login" {
	// 			logger.Info("✓ Security verification completed successfully!")
	// 			return nil
	// 		}
	// 		logger.Info(fmt.Sprintf("Still waiting... (%d/30 checks)", i+1))
	// 	}
	// 	return errors.New("Security verification timeout - please try again")
	// }

	// Final check - are we logged in now?
	currentURL = page.MustInfo().URL
	logger.Info("Final URL check: " + currentURL)

	// LinkedIn home page URL should contain "/feed" or similar indicators
	if currentURL != "https://www.linkedin.com/login" {
		logger.Info("Login Successful - Redirected to home page")
		return nil
	}

	return errors.New("login failed - still on login page")
}

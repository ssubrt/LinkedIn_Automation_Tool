package automation

/*
	errors- used to return meaningful errors
	rod- browser interaction library
	stealth - humanlike behaviour functions
	logger - cleans logs instead of print

*/
import (
	"errors"

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
	emailInput, err := page.Element(`input#username`)
	if err != nil {
		return errors.New("email input not found")
	}

	//pause for random time to mimic human behaviour and type email id like a human
	stealth.RandomDelay(800, 1500)
	emailInput.MustInput(email)

	//Locate password input field  and tell user if the password input field is empty
	logger.Info("Locating password input field")
	passwordInput, err := page.Element(`input#password`)
	if err != nil {
		return errors.New("password input not found")
	}

	//pause for random time to mimic human behaviour and type password like a human
	stealth.RandomDelay(800, 1500)
	stealth.TypeLikeHuman(passwordInput, password)

	//Locate and click on the Sign in button to submit the credentials
	logger.Info("Locating and clicking on sign in button")
	loginBtn, err := page.Element(`button[type="submit"]`)
	if err != nil {
		return errors.New("Login Button not found")
	}

	//puase for random time to show human behaviour and  click on the login button
	stealth.RandomDelay(1000, 2000)
	loginBtn.MustClick()

	//Wait for linkedln page to load after the login btn is clicked
	stealth.RandomDelay(3000, 5000)
	page.MustWaitLoad()

	// Check if login was successful by checking page URL
	logger.Info("Checking if login was successful")

	// Wait a bit for page to fully stabilize after login
	stealth.RandomDelay(2000, 3000)

	// Get the current URL to verify we're on the home page
	currentURL := page.MustInfo().URL
	logger.Info("Current page URL: " + currentURL)

	// LinkedIn home page URL should contain "/feed" or similar indicators
	if currentURL != "https://www.linkedin.com/login" {
		logger.Info("Login Successful - Redirected to home page")
		return nil
	}

	return errors.New("login failed - still on login page")
}

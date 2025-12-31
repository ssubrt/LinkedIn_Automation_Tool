package utils

// Constants for LinkedIn automation
const (
	// LinkedIn URLs
	LinkedInBaseURL     = "https://www.linkedin.com"
	LinkedInLoginURL    = "https://www.linkedin.com/login"
	LinkedInFeedURL     = "https://www.linkedin.com/feed/"
	LinkedInSearchURL   = "https://www.linkedin.com/search/results/people/"
	LinkedInProfileBase = "https://www.linkedin.com/in/"

	// Delay ranges (milliseconds)
	MinLoginDelay  = 800
	MaxLoginDelay  = 1500
	MinScrollDelay = 800
	MaxScrollDelay = 1500
	MinMouseDelay  = 300
	MaxMouseDelay  = 800

	// Mouse movement ranges
	MinMouseX = 100
	MaxMouseX = 800
	MinMouseY = 100
	MaxMouseY = 600

	// Scroll ranges
	MinScrollDist = 200
	MaxScrollDist = 600

	// Retry settings
	MaxRetries   = 3
	RetryDelayMS = 2000

	// Timeouts
	PageLoadTimeout = 30
	LoginTimeout    = 60
	DefaultTimeout  = 30
)

// Error messages
const (
	ErrorInvalidEmail    = "invalid email format"
	ErrorInvalidPassword = "invalid password"
	ErrorLoginFailed     = "login failed"
	ErrorBrowserCrash    = "browser crashed"
	ErrorTimeout         = "operation timed out"
	ErrorNetwork         = "network error"
)

// Action types
const (
	ActionLogin     = "login"
	ActionScroll    = "scroll"
	ActionMouseMove = "mouse_move"
	ActionType      = "type"
	ActionClick     = "click"
	ActionWait      = "wait"
)

// Browser settings
const (
	DefaultBrowserTimeout = 30
	ChromeUserAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

// Stealth modes
const (
	StealthModeOff      = "off"
	StealthModeBasic    = "basic"
	StealthModeAdvanced = "advanced"
	StealthModeMaximum  = "maximum"
)

// Log levels
const (
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelError = "ERROR"
	LogLevelFatal = "FATAL"
)

// LinkedIn Location URN Codes (geoUrn parameter)
// These are LinkedIn's internal IDs for geographic locations
var LinkedInLocations = map[string]string{
	// United States - Major Cities
	"San Francisco Bay Area": "90000084",
	"New York City Area":     "90000070",
	"Los Angeles":            "90000071",
	"Chicago":                "90000074",
	"Boston":                 "90000075",
	"Seattle":                "90000076",
	"Austin":                 "90000073",
	"Denver":                 "90000077",
	"Washington DC":          "90000078",
	"Atlanta":                "90000079",
	"Dallas":                 "90000080",
	"Miami":                  "90000081",
	"Philadelphia":           "90000082",
	"Phoenix":                "90000083",
	"San Diego":              "90000085",

	// United States - States
	"California":    "102095887",
	"New York":      "105080838",
	"Texas":         "102748797",
	"Florida":       "104022003",
	"Illinois":      "104677192",
	"Massachusetts": "104842724",
	"Washington":    "103977809",
	"Colorado":      "104831318",

	// United States - Country
	"United States": "103644278",

	// International - Countries
	"United Kingdom": "101165590",
	"Canada":         "101174742",
	"Germany":        "101282230",
	"France":         "105015875",
	"India":          "102713980",
	"Australia":      "101452733",
	"Netherlands":    "102890719",
	"Singapore":      "102454443",
	"Brazil":         "106057199",
	"Japan":          "101355337",
	"China":          "102890883",
	"Spain":          "105646813",
	"Italy":          "103350119",
	"Mexico":         "103323778",

	// International - Major Cities
	"London":         "90009496",
	"Toronto":        "90009496",
	"Berlin":         "106967730",
	"Paris":          "105015875",
	"Sydney":         "104769905",
	"Bangalore":      "105214831",
	"Amsterdam":      "100561920",
	"Singapore City": "102454443",
	"Tokyo":          "104738515",
	"Hong Kong":      "102279293",
	"Dubai":          "104305776",
	"Munich":         "106693272",
	"Barcelona":      "100994330",
	"Madrid":         "103924744",
}

// Search result selectors
// ⚠️  WARNING: LinkedIn changes these selectors frequently (every 3-6 months)
// If search returns 0 results, check the browser inspector and update these:
// 1. Open LinkedIn search in browser
// 2. Right-click on a profile card → Inspect
// 3. Find the updated class names for profile containers
// 4. Update the constants below
// Last verified: December 2025
const (
	SearchResultContainerSelector = ".reusable-search__result-container" // Alternative: .search-results-container
	SearchResultItemSelector      = ".entity-result"                     // Alternative: .search-result__info
	SearchResultTitleSelector     = ".entity-result__title-text a"       // Alternative: .app-aware-link
	SearchResultSubtitleSelector  = ".entity-result__primary-subtitle"   // Alternative: .entity-result__subtitle
	SearchResultSecondarySelector = ".entity-result__secondary-subtitle" // Alternative: .entity-result__summary
	SearchResultLinkSelector      = "a.app-aware-link"                   // Alternative: a[href*='/in/']
	PaginationNextButtonSelector  = ".artdeco-pagination__button--next"  // Alternative: button[aria-label='Next']
	PaginationDisabledClass       = "artdeco-button--disabled"           // Check for 'disabled' attribute too
)

// Search constraints
const (
	MaxSearchResultsPerPage = 10
	MaxPaginationPages      = 100
	SearchDelaySeconds      = 2
)

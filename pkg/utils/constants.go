package utils

// Constants for LinkedIn automation
const (
	// LinkedIn URLs
	LinkedInBaseURL  = "https://www.linkedin.com"
	LinkedInLoginURL = "https://www.linkedin.com/login"
	LinkedInFeedURL  = "https://www.linkedin.com/feed/"

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
	MaxRetries      = 3
	RetryDelayMS    = 2000

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
	StealthModeOff       = "off"
	StealthModeBasic     = "basic"
	StealthModeAdvanced  = "advanced"
	StealthModeMaximum   = "maximum"
)

// Log levels
const (
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelError = "ERROR"
	LogLevelFatal = "FATAL"
)

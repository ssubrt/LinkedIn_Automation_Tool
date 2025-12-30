package models

import "time"

// User represents a LinkedIn user account
type User struct {
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"` // Don't expose in JSON
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Profile   Profile   `json:"profile,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Profile contains LinkedIn profile information
type Profile struct {
	Headline        string `json:"headline"`
	Bio             string `json:"bio"`
	Location        string `json:"location"`
	ProfilePicURL   string `json:"profile_pic_url"`
	ConnectionCount int    `json:"connection_count"`
}

// LoginAttempt tracks login activity
type LoginAttempt struct {
	Email       string    `json:"email"`
	Timestamp   time.Time `json:"timestamp"`
	Success     bool      `json:"success"`
	ErrorReason string    `json:"error_reason,omitempty"`
	IPAddress   string    `json:"ip_address,omitempty"`
}

// AutomationConfig contains automation settings
type AutomationConfig struct {
	HeadlessBrowser bool   `json:"headless_browser"`
	StealthMode     bool   `json:"stealth_mode"`
	MaxRetries      int    `json:"max_retries"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
	EnableLogging   bool   `json:"enable_logging"`
	LogFilePath     string `json:"log_file_path"`
	RandomDelayMin  int    `json:"random_delay_min"`
	RandomDelayMax  int    `json:"random_delay_max"`
}

// AutomationResult represents the outcome of automation
type AutomationResult struct {
	Success   bool      `json:"success"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int64     `json:"duration_ms"`
	Message   string    `json:"message"`
	Error     string    `json:"error,omitempty"`
	Actions   []Action  `json:"actions"`
}

// Action represents a single automation action
type Action struct {
	Type      string    `json:"type"` // login, scroll, move, type, click
	Target    string    `json:"target"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
	Details   string    `json:"details,omitempty"`
}

// MouseMovement represents a mouse movement action
type MouseMovement struct {
	X        float64   `json:"x"`
	Y        float64   `json:"y"`
	Duration int       `json:"duration_ms"`
	Time     time.Time `json:"timestamp"`
}

// ScrollAction represents a scroll action
type ScrollAction struct {
	Direction string    `json:"direction"` // up, down
	Distance  int       `json:"distance"`
	Duration  int       `json:"duration_ms"`
	Time      time.Time `json:"timestamp"`
}

// StealthTechnique represents anti-detection measures
type StealthTechnique struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

// Report represents automation report
type Report struct {
	ID             string             `json:"id"`
	Timestamp      time.Time          `json:"timestamp"`
	User           User               `json:"user"`
	Result         AutomationResult   `json:"result"`
	MouseMovements []MouseMovement    `json:"mouse_movements"`
	ScrollActions  []ScrollAction     `json:"scroll_actions"`
	StealthMethods []StealthTechnique `json:"stealth_methods"`
}

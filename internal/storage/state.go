package storage

import (
	"encoding/json"
	"os"
	"time"
)

// AppState represents the current state of the LinkedIn automation application.
// It tracks whether a login was attempted, the timestamp of the last run, and session validity.
type AppState struct {
	// LoginAttempted indicates whether a login attempt was made
	LoginAttempted bool `json:"login_attempted"`
	// LastRun stores the timestamp of when the automation last ran
	LastRun time.Time `json:"last_run"`
	// SessionValid indicates if the browser session is still active
	SessionValid bool `json:"session_valid"`
	// LastLoginTime stores when the last successful login occurred
	LastLoginTime time.Time `json:"last_login_time"`
	// BrowserDataDir stores the path to the persistent browser data directory
	BrowserDataDir string `json:"browser_data_dir"`
}

const stateFilePath = "data/state.json"

// SaveState saves the current application state to a JSON file.
// It creates or overwrites the data/state.json file with the current timestamp and login status.
// Returns an error if file creation or encoding fails.
func SaveState(sessionValid bool) error {
	// Load existing state to preserve certain fields
	existingState, _ := LoadState()

	// Create an AppState struct with current timestamp and login status
	state := AppState{
		LoginAttempted: true,
		LastRun:        time.Now(),
		SessionValid:   sessionValid,
		LastLoginTime:  time.Now(),
		BrowserDataDir: "./browser_data",
	}

	// Preserve last login time if session was already valid
	if existingState != nil && existingState.SessionValid {
		state.LastLoginTime = existingState.LastLoginTime
	}

	// Ensure the data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	// Create or truncate the state file at the specified path
	file, err := os.Create(stateFilePath)
	if err != nil {
		return err
	}

	// Ensure the file is closed when the function returns
	defer file.Close()

	// Encode the state struct to JSON with indentation for readability
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(state)
}

// LoadState loads the application state from the JSON file.
// Returns the AppState struct if the file exists, or nil if not found.
// Returns an error if file reading or decoding fails.
func LoadState() (*AppState, error) {
	// Check if state file exists
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return nil, nil // File doesn't exist, return nil (not an error)
	}

	// Open the state file
	file, err := os.Open(stateFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the JSON file into AppState struct
	var state AppState
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&state); err != nil {
		return nil, err
	}

	return &state, nil
}

// IsSessionValid checks if the saved session is still valid (less than 7 days old)
func IsSessionValid(state *AppState) bool {
	if state == nil || !state.SessionValid {
		return false
	}

	// Session is valid if last login was within 7 days
	sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)
	return state.LastLoginTime.After(sevenDaysAgo)
}

// InvalidateSession marks the current session as invalid
func InvalidateSession() error {
	state, err := LoadState()
	if err != nil || state == nil {
		// If no state exists, create a new one
		state = &AppState{
			LoginAttempted: false,
			LastRun:        time.Now(),
			SessionValid:   false,
			BrowserDataDir: "./browser_data",
		}
	} else {
		state.SessionValid = false
		state.LastRun = time.Now()
	}

	// Save the updated state
	return SaveState(false)
}

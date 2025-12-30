package storage

import (
	"encoding/json"
	"os"
	"testing"
)

// TestSaveState verifies state file is created correctly
func TestSaveState(t *testing.T) {
	// Save state
	err := SaveState(true)
	if err != nil {
		t.Errorf("SaveState failed: %v", err)
	}

	// Verify file exists
	_, err = os.Stat("data/state.json")
	if os.IsNotExist(err) {
		t.Error("State file was not created")
	}
}

// TestAppStateStructure verifies AppState has correct fields
func TestAppStateStructure(t *testing.T) {
	state := AppState{
		LoginAttempted: true,
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(state)
	if err != nil {
		t.Errorf("Failed to marshal AppState: %v", err)
	}

	// Verify JSON contains expected fields
	var jsonMap map[string]interface{}
	json.Unmarshal(jsonBytes, &jsonMap)

	if _, ok := jsonMap["login_attempted"]; !ok {
		t.Error("AppState missing login_attempted field")
	}
	if _, ok := jsonMap["last_run"]; !ok {
		t.Error("AppState missing last_run field")
	}
}

// TestAppStateJSON verifies JSON marshaling/unmarshaling
func TestAppStateJSON(t *testing.T) {
	original := AppState{
		LoginAttempted: true,
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(original)
	if err != nil {
		t.Errorf("Failed to marshal: %v", err)
	}

	// Unmarshal back
	var restored AppState
	err = json.Unmarshal(jsonBytes, &restored)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	// Verify values match
	if original.LoginAttempted != restored.LoginAttempted {
		t.Errorf("LoginAttempted mismatch: %v != %v", original.LoginAttempted, restored.LoginAttempted)
	}
}

// TestStateFileCreation ensures data directory exists
func TestStateFileCreation(t *testing.T) {
	// Create data directory if it doesn't exist
	os.MkdirAll("data", 0755)

	err := SaveState(true)
	if err != nil {
		t.Errorf("Failed to save state: %v", err)
	}

	// Verify file is readable
	_, err = os.ReadFile("data/state.json")
	if err != nil {
		t.Errorf("Failed to read state file: %v", err)
	}
}

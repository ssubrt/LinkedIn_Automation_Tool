package storage

import (
	"os"
	"testing"
	"time"
)

func TestInitDB(t *testing.T) {
	// Use a temporary database for testing
	testDBPath := "./test_linkedin.db"
	defer os.Remove(testDBPath) // Clean up after test

	db, err := InitDB(testDBPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Verify database connection
	if db.conn == nil {
		t.Error("Database connection is nil")
	}
}

func TestSaveAndGetProfile(t *testing.T) {
	testDBPath := "./test_linkedin.db"
	defer os.Remove(testDBPath)

	db, err := InitDB(testDBPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create test profile
	profile := Profile{
		ID:         "test-profile-123",
		Name:       "John Doe",
		Title:      "Software Engineer",
		Company:    "Tech Corp",
		Location:   "San Francisco",
		ProfileURL: "https://linkedin.com/in/johndoe",
		VisitedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}

	// Save profile
	err = db.SaveProfile(profile)
	if err != nil {
		t.Errorf("Failed to save profile: %v", err)
	}

	// Retrieve profile
	retrieved, err := db.GetProfile("test-profile-123")
	if err != nil {
		t.Errorf("Failed to get profile: %v", err)
	}

	// Verify data
	if retrieved.Name != profile.Name {
		t.Errorf("Name mismatch: expected %s, got %s", profile.Name, retrieved.Name)
	}
	if retrieved.Title != profile.Title {
		t.Errorf("Title mismatch: expected %s, got %s", profile.Title, retrieved.Title)
	}
}

func TestIsDuplicateProfile(t *testing.T) {
	testDBPath := "./test_linkedin.db"
	defer os.Remove(testDBPath)

	db, err := InitDB(testDBPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create and save profile
	profile := Profile{
		ID:         "duplicate-test-123",
		Name:       "Jane Smith",
		Title:      "Product Manager",
		Company:    "Startup Inc",
		Location:   "New York",
		ProfileURL: "https://linkedin.com/in/janesmith",
		VisitedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}

	err = db.SaveProfile(profile)
	if err != nil {
		t.Errorf("Failed to save profile: %v", err)
	}

	// Check if it's a duplicate (within 30 days)
	isDuplicate, err := db.IsDuplicateProfile("duplicate-test-123", 30)
	if err != nil {
		t.Errorf("Failed to check duplicate: %v", err)
	}

	if !isDuplicate {
		t.Error("Profile should be detected as duplicate")
	}

	// Check non-existent profile
	isDuplicate, err = db.IsDuplicateProfile("non-existent-id", 30)
	if err != nil {
		t.Errorf("Failed to check duplicate: %v", err)
	}

	if isDuplicate {
		t.Error("Non-existent profile should not be duplicate")
	}
}

func TestSaveConnectionRequest(t *testing.T) {
	testDBPath := "./test_linkedin.db"
	defer os.Remove(testDBPath)

	db, err := InitDB(testDBPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Save connection request
	req := ConnectionRequest{
		ProfileID: "test-profile-123",
		SentAt:    time.Now(),
		NoteUsed:  "Hi, I'd like to connect!",
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	err = db.SaveConnectionRequest(req)
	if err != nil {
		t.Errorf("Failed to save connection request: %v", err)
	}

	// Check if request exists
	hasSent, err := db.HasSentConnectionRequest("test-profile-123")
	if err != nil {
		t.Errorf("Failed to check connection request: %v", err)
	}

	if !hasSent {
		t.Error("Connection request should exist")
	}
}

func TestRateLimits(t *testing.T) {
	testDBPath := "./test_linkedin.db"
	defer os.Remove(testDBPath)

	db, err := InitDB(testDBPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Get today's rate limit (should create new record)
	limit, err := db.GetTodayRateLimit()
	if err != nil {
		t.Errorf("Failed to get rate limit: %v", err)
	}

	if limit.ConnectionCount != 0 {
		t.Errorf("Initial connection count should be 0, got %d", limit.ConnectionCount)
	}

	// Increment connection count
	err = db.IncrementConnectionCount()
	if err != nil {
		t.Errorf("Failed to increment connection count: %v", err)
	}

	// Verify increment
	limit, err = db.GetTodayRateLimit()
	if err != nil {
		t.Errorf("Failed to get rate limit: %v", err)
	}

	if limit.ConnectionCount != 1 {
		t.Errorf("Connection count should be 1, got %d", limit.ConnectionCount)
	}

	// Increment message count
	err = db.IncrementMessageCount()
	if err != nil {
		t.Errorf("Failed to increment message count: %v", err)
	}

	// Verify both counts
	limit, err = db.GetTodayRateLimit()
	if err != nil {
		t.Errorf("Failed to get rate limit: %v", err)
	}

	if limit.ConnectionCount != 1 || limit.MessageCount != 1 {
		t.Errorf("Expected counts 1,1 got %d,%d", limit.ConnectionCount, limit.MessageCount)
	}
}

func TestSaveAndRetrieveMessage(t *testing.T) {
	testDBPath := "./test_linkedin.db"
	defer os.Remove(testDBPath)

	db, err := InitDB(testDBPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Save message
	msg := Message{
		ConnectionID:   "connection-123",
		TemplateName:   "welcome",
		MessageContent: "Hi! Thanks for connecting.",
		SentAt:         time.Now(),
		CreatedAt:      time.Now(),
	}

	err = db.SaveMessage(msg)
	if err != nil {
		t.Errorf("Failed to save message: %v", err)
	}

	// Check if message exists
	hasSent, err := db.HasSentMessage("connection-123", "welcome")
	if err != nil {
		t.Errorf("Failed to check message: %v", err)
	}

	if !hasSent {
		t.Error("Message should exist")
	}

	// Get message history
	history, err := db.GetMessageHistory("connection-123")
	if err != nil {
		t.Errorf("Failed to get message history: %v", err)
	}

	if len(history) != 1 {
		t.Errorf("Expected 1 message in history, got %d", len(history))
	}

	if history[0].TemplateName != "welcome" {
		t.Errorf("Expected template 'welcome', got '%s'", history[0].TemplateName)
	}
}

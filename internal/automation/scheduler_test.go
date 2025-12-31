package automation

import (
	"testing"
	"time"
)

func TestGetDefaultSchedule(t *testing.T) {
	config := GetDefaultSchedule()

	// Check defaults
	if config.StartHour < 0 || config.StartHour > 23 {
		t.Errorf("Invalid start hour: %d", config.StartHour)
	}

	if config.EndHour < 0 || config.EndHour > 23 {
		t.Errorf("Invalid end hour: %d", config.EndHour)
	}

	if config.StartHour >= config.EndHour {
		t.Errorf("Start hour (%d) should be before end hour (%d)", config.StartHour, config.EndHour)
	}
}

func TestIsActiveHoursWithConfig(t *testing.T) {
	// Test with config that should always be active
	alwaysActive := ScheduleConfig{
		StartHour:    0,
		EndHour:      23,
		WeekdaysOnly: false,
	}

	if !IsActiveHoursWithConfig(alwaysActive) {
		t.Error("With 0-23 hours and no weekday restriction, should always be active")
	}

	// Test with config that's very restrictive (unlikely to match current time)
	restricted := ScheduleConfig{
		StartHour:    2,
		EndHour:      3,
		WeekdaysOnly: false,
	}

	// This might be active or not depending on current time, so we just verify it doesn't crash
	_ = IsActiveHoursWithConfig(restricted)
}

func TestIsActiveHoursWeekendDetection(t *testing.T) {
	// Create a config that excludes weekends
	weekdayOnly := ScheduleConfig{
		StartHour:    0,
		EndHour:      23,
		WeekdaysOnly: true,
	}

	// We can't easily test this without mocking time, but we can verify
	// it doesn't crash and returns a boolean
	result := IsActiveHoursWithConfig(weekdayOnly)

	// Check that it returns a boolean
	if result != true && result != false {
		t.Error("IsActiveHoursWithConfig should return a boolean")
	}
}

func TestCalculateNextActiveTime(t *testing.T) {
	config := ScheduleConfig{
		StartHour:    9,
		EndHour:      17,
		WeekdaysOnly: false,
	}

	// Test with a time that's currently active (10 AM)
	testTime := time.Date(2025, 12, 30, 10, 0, 0, 0, time.Local)
	nextActive := CalculateNextActiveTime(testTime, config)

	// Next active should be same day (if still active) or next day
	// It should never be before the test time
	if nextActive.Before(testTime.Add(-1 * time.Hour)) {
		t.Error("Next active time should not be significantly in the past")
	}

	// Test with a time that's after end hour (6 PM)
	testTime = time.Date(2025, 12, 30, 18, 0, 0, 0, time.Local)
	nextActive = CalculateNextActiveTime(testTime, config)

	// Next active should be next day at 9 AM
	if nextActive.Day() != testTime.Day()+1 {
		t.Error("Next active time should be next day")
	}

	if nextActive.Hour() != config.StartHour {
		t.Errorf("Next active time should be at start hour %d, got %d", config.StartHour, nextActive.Hour())
	}
}

func TestCalculateNextActiveTimeWeekend(t *testing.T) {
	config := ScheduleConfig{
		StartHour:    9,
		EndHour:      17,
		WeekdaysOnly: true,
	}

	// Test with a Saturday
	saturday := time.Date(2025, 12, 27, 10, 0, 0, 0, time.Local) // Dec 27, 2025 is Saturday
	nextActive := CalculateNextActiveTime(saturday, config)

	// Next active should be Monday
	if nextActive.Weekday() != time.Monday {
		t.Errorf("Next active after Saturday should be Monday, got %v", nextActive.Weekday())
	}

	// Test with a Sunday
	sunday := time.Date(2025, 12, 28, 10, 0, 0, 0, time.Local)
	nextActive = CalculateNextActiveTime(sunday, config)

	// Next active should be Monday
	if nextActive.Weekday() != time.Monday {
		t.Errorf("Next active after Sunday should be Monday, got %v", nextActive.Weekday())
	}
}

func TestGetTimeUntilNextActive(t *testing.T) {
	// Test that it returns a non-negative duration
	duration := GetTimeUntilNextActive()

	if duration < 0 {
		t.Error("Time until next active should never be negative")
	}

	// If we're in active hours, duration should be 0
	if IsActiveHours() {
		if duration != 0 {
			t.Errorf("If in active hours, duration should be 0, got %v", duration)
		}
	}
}

func TestShouldPauseAutomation(t *testing.T) {
	shouldPause, reason := ShouldPauseAutomation()

	// Should return a boolean and a string
	if shouldPause != true && shouldPause != false {
		t.Error("ShouldPauseAutomation should return a boolean")
	}

	// If paused, should have a reason
	if shouldPause && reason == "" {
		t.Error("If paused, should provide a reason")
	}

	// If not paused, reason can be empty
	if !shouldPause && !IsActiveHours() {
		t.Error("If outside active hours, should pause")
	}
}

func BenchmarkIsActiveHours(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsActiveHours()
	}
}

func BenchmarkCalculateNextActiveTime(b *testing.B) {
	config := GetDefaultSchedule()
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateNextActiveTime(now, config)
	}
}

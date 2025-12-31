package automation

import (
	"os"
	"strconv"
	"time"

	"linkedin-automation/internal/logger"
)

// ScheduleConfig holds configuration for activity scheduling
type ScheduleConfig struct {
	StartHour    int  // Business hours start (default: 9 AM)
	EndHour      int  // Business hours end (default: 5 PM)
	WeekdaysOnly bool // Only operate on weekdays (Monday-Friday)
}

// GetDefaultSchedule returns the default scheduling configuration
func GetDefaultSchedule() ScheduleConfig {
	// Try to get from environment variables
	startHour := 9
	endHour := 17
	weekdaysOnly := true

	if envStart := os.Getenv("ACTIVE_HOURS_START"); envStart != "" {
		if h, err := strconv.Atoi(envStart); err == nil && h >= 0 && h < 24 {
			startHour = h
		}
	}

	if envEnd := os.Getenv("ACTIVE_HOURS_END"); envEnd != "" {
		if h, err := strconv.Atoi(envEnd); err == nil && h >= 0 && h < 24 {
			endHour = h
		}
	}

	if envWeekdays := os.Getenv("WEEKDAYS_ONLY"); envWeekdays != "" {
		weekdaysOnly = envWeekdays == "true"
	}

	return ScheduleConfig{
		StartHour:    startHour,
		EndHour:      endHour,
		WeekdaysOnly: weekdaysOnly,
	}
}

// IsActiveHours checks if the current time is within business hours
// Returns true if automation should run, false otherwise
func IsActiveHours() bool {
	return IsActiveHoursWithConfig(GetDefaultSchedule())
}

// IsActiveHoursWithConfig checks if the current time is within configured hours
func IsActiveHoursWithConfig(config ScheduleConfig) bool {
	now := time.Now()

	// Check if it's a weekday (Monday = 1, Sunday = 0)
	if config.WeekdaysOnly {
		weekday := now.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			logger.Debug("Outside active hours: Weekend detected")
			return false
		}
	}

	// Check if it's within business hours
	currentHour := now.Hour()
	if currentHour < config.StartHour || currentHour >= config.EndHour {
		logger.Debug("Outside active hours: Current hour " + strconv.Itoa(currentHour) +
			" not in range " + strconv.Itoa(config.StartHour) + "-" + strconv.Itoa(config.EndHour))
		return false
	}

	return true
}

// WaitForActiveHours blocks execution until we're in active hours
// Returns immediately if already in active hours
func WaitForActiveHours() {
	WaitForActiveHoursWithConfig(GetDefaultSchedule())
}

// WaitForActiveHoursWithConfig blocks until configured active hours
func WaitForActiveHoursWithConfig(config ScheduleConfig) {
	if IsActiveHoursWithConfig(config) {
		return
	}

	now := time.Now()

	// Calculate next active time
	nextActive := CalculateNextActiveTime(now, config)

	waitDuration := nextActive.Sub(now)
	logger.Info("Outside active hours. Waiting until " + nextActive.Format("2006-01-02 15:04:05") +
		" (" + waitDuration.String() + ")")

	time.Sleep(waitDuration)

	logger.Info("Active hours resumed")
}

// CalculateNextActiveTime calculates the next time when automation should run
func CalculateNextActiveTime(current time.Time, config ScheduleConfig) time.Time {
	// Start with today at the start hour
	nextActive := time.Date(
		current.Year(), current.Month(), current.Day(),
		config.StartHour, 0, 0, 0, current.Location(),
	)

	// If we're already past start hour today, move to tomorrow
	if current.Hour() >= config.EndHour {
		nextActive = nextActive.Add(24 * time.Hour)
	}

	// Skip weekends if configured
	if config.WeekdaysOnly {
		for {
			weekday := nextActive.Weekday()
			if weekday == time.Saturday {
				// Skip to Monday
				nextActive = nextActive.Add(48 * time.Hour)
			} else if weekday == time.Sunday {
				// Skip to Monday
				nextActive = nextActive.Add(24 * time.Hour)
			} else {
				break
			}
		}
	}

	return nextActive
}

// GetTimeUntilNextActive returns the duration until next active hours
func GetTimeUntilNextActive() time.Duration {
	return GetTimeUntilNextActiveWithConfig(GetDefaultSchedule())
}

// GetTimeUntilNextActiveWithConfig returns duration until next active hours
func GetTimeUntilNextActiveWithConfig(config ScheduleConfig) time.Duration {
	if IsActiveHoursWithConfig(config) {
		return 0
	}

	now := time.Now()
	nextActive := CalculateNextActiveTime(now, config)
	return nextActive.Sub(now)
}

// ShouldPauseAutomation checks if automation should pause
// This can be extended to check for other conditions like rate limits
func ShouldPauseAutomation() (bool, string) {
	if !IsActiveHours() {
		return true, "Outside active hours"
	}

	// Can add more conditions here:
	// - Check if rate limits exceeded
	// - Check if maintenance window
	// - Check if error rate too high

	return false, ""
}

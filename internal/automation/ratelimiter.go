package automation

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"linkedin-automation/internal/logger"
	"linkedin-automation/internal/storage"
)

// TaskType represents different types of automation tasks
type TaskType string

const (
	TaskConnection TaskType = "connection"
	TaskMessage    TaskType = "message"
	TaskSearch     TaskType = "search"
)

// RateLimitConfig holds rate limit settings
type RateLimitConfig struct {
	MaxConnectionsPerDay   int
	MaxMessagesPerDay      int
	MaxSearchesPerDay      int
	CooldownBetweenActions time.Duration // Cooldown between individual actions
}

// RateLimitError represents a rate limit exceeded error
type RateLimitError struct {
	TaskType  TaskType
	Current   int
	Limit     int
	ResetTime time.Time
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("Rate limit exceeded for %s: %d/%d (resets at %s)",
		e.TaskType, e.Current, e.Limit, e.ResetTime.Format("15:04:05"))
}

// GetDefaultRateLimitConfig returns default rate limits from env or constants
func GetDefaultRateLimitConfig() RateLimitConfig {
	config := RateLimitConfig{
		MaxConnectionsPerDay:   14,               // Safe default: ~100/week
		MaxMessagesPerDay:      50,               // LinkedIn's typical limit
		MaxSearchesPerDay:      100,              // Conservative search limit
		CooldownBetweenActions: 30 * time.Second, // 30s cooldown between actions
	}

	// Override from environment variables
	if envConn := os.Getenv("MAX_CONNECTIONS_PER_DAY"); envConn != "" {
		if val, err := strconv.Atoi(envConn); err == nil && val > 0 {
			config.MaxConnectionsPerDay = val
		}
	}

	if envMsg := os.Getenv("MAX_MESSAGES_PER_DAY"); envMsg != "" {
		if val, err := strconv.Atoi(envMsg); err == nil && val > 0 {
			config.MaxMessagesPerDay = val
		}
	}

	if envSearch := os.Getenv("MAX_SEARCHES_PER_DAY"); envSearch != "" {
		if val, err := strconv.Atoi(envSearch); err == nil && val > 0 {
			config.MaxSearchesPerDay = val
		}
	}

	if envCooldown := os.Getenv("COOLDOWN_SECONDS"); envCooldown != "" {
		if val, err := strconv.Atoi(envCooldown); err == nil && val > 0 {
			config.CooldownBetweenActions = time.Duration(val) * time.Second
		}
	}

	return config
}

// RateLimiter manages rate limiting for automation tasks
type RateLimiter struct {
	db             *storage.Database
	config         RateLimitConfig
	lastActionTime time.Time
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(db *storage.Database) *RateLimiter {
	return &RateLimiter{
		db:             db,
		config:         GetDefaultRateLimitConfig(),
		lastActionTime: time.Now().Add(-1 * time.Hour), // Allow immediate first action
	}
}

// NewRateLimiterWithConfig creates a rate limiter with custom config
func NewRateLimiterWithConfig(db *storage.Database, config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		db:             db,
		config:         config,
		lastActionTime: time.Now().Add(-1 * time.Hour),
	}
}

// CheckDailyLimit checks if a task type has exceeded its daily limit
// Returns error if limit exceeded, nil otherwise
func (rl *RateLimiter) CheckDailyLimit(taskType TaskType) error {
	// Get today's rate limit from database
	limit, err := rl.db.GetTodayRateLimit()
	if err != nil {
		return fmt.Errorf("failed to get rate limit: %w", err)
	}

	// Check limit based on task type
	switch taskType {
	case TaskConnection:
		if limit.ConnectionCount >= rl.config.MaxConnectionsPerDay {
			return &RateLimitError{
				TaskType:  TaskConnection,
				Current:   limit.ConnectionCount,
				Limit:     rl.config.MaxConnectionsPerDay,
				ResetTime: rl.getNextMidnight(),
			}
		}
	case TaskMessage:
		if limit.MessageCount >= rl.config.MaxMessagesPerDay {
			return &RateLimitError{
				TaskType:  TaskMessage,
				Current:   limit.MessageCount,
				Limit:     rl.config.MaxMessagesPerDay,
				ResetTime: rl.getNextMidnight(),
			}
		}
	case TaskSearch:
		if limit.SearchCount >= rl.config.MaxSearchesPerDay {
			return &RateLimitError{
				TaskType:  TaskSearch,
				Current:   limit.SearchCount,
				Limit:     rl.config.MaxSearchesPerDay,
				ResetTime: rl.getNextMidnight(),
			}
		}
	default:
		return fmt.Errorf("unknown task type: %s", taskType)
	}

	return nil
}

// ApplyCooldown waits for the cooldown period since last action
func (rl *RateLimiter) ApplyCooldown() {
	timeSinceLastAction := time.Since(rl.lastActionTime)

	if timeSinceLastAction < rl.config.CooldownBetweenActions {
		waitTime := rl.config.CooldownBetweenActions - timeSinceLastAction
		logger.Info(fmt.Sprintf("Applying cooldown: waiting %.1f seconds", waitTime.Seconds()))
		time.Sleep(waitTime)
	}

	rl.lastActionTime = time.Now()
}

// RecordAction records that an action was performed and increments the counter
func (rl *RateLimiter) RecordAction(taskType TaskType) error {
	// Apply cooldown before action
	rl.ApplyCooldown()

	// Increment the counter in database
	var err error
	switch taskType {
	case TaskConnection:
		err = rl.db.IncrementConnectionCount()
	case TaskMessage:
		err = rl.db.IncrementMessageCount()
	case TaskSearch:
		err = rl.db.IncrementSearchCount()
	default:
		return fmt.Errorf("unknown task type: %s", taskType)
	}

	if err != nil {
		return fmt.Errorf("failed to record action: %w", err)
	}

	return nil
}

// GetRemainingQuota returns how many actions are remaining for a task type
func (rl *RateLimiter) GetRemainingQuota(taskType TaskType) (int, error) {
	limit, err := rl.db.GetTodayRateLimit()
	if err != nil {
		return 0, err
	}

	switch taskType {
	case TaskConnection:
		return rl.config.MaxConnectionsPerDay - limit.ConnectionCount, nil
	case TaskMessage:
		return rl.config.MaxMessagesPerDay - limit.MessageCount, nil
	case TaskSearch:
		return rl.config.MaxSearchesPerDay - limit.SearchCount, nil
	default:
		return 0, fmt.Errorf("unknown task type: %s", taskType)
	}
}

// GetUsagePercentage returns the percentage of daily quota used
func (rl *RateLimiter) GetUsagePercentage(taskType TaskType) (float64, error) {
	limit, err := rl.db.GetTodayRateLimit()
	if err != nil {
		return 0, err
	}

	var current, max int
	switch taskType {
	case TaskConnection:
		current = limit.ConnectionCount
		max = rl.config.MaxConnectionsPerDay
	case TaskMessage:
		current = limit.MessageCount
		max = rl.config.MaxMessagesPerDay
	case TaskSearch:
		current = limit.SearchCount
		max = rl.config.MaxSearchesPerDay
	default:
		return 0, fmt.Errorf("unknown task type: %s", taskType)
	}

	if max == 0 {
		return 0, nil
	}

	return float64(current) / float64(max) * 100, nil
}

// ShouldWarnAboutLimit checks if we're approaching the limit (80% threshold)
func (rl *RateLimiter) ShouldWarnAboutLimit(taskType TaskType) (bool, error) {
	percentage, err := rl.GetUsagePercentage(taskType)
	if err != nil {
		return false, err
	}

	return percentage >= 80.0, nil
}

// getNextMidnight returns the time of the next midnight (when limits reset)
func (rl *RateLimiter) getNextMidnight() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
}

// GetDailyStats returns a summary of today's rate limit usage
func (rl *RateLimiter) GetDailyStats() (string, error) {
	limit, err := rl.db.GetTodayRateLimit()
	if err != nil {
		return "", err
	}

	connPercent, _ := rl.GetUsagePercentage(TaskConnection)
	msgPercent, _ := rl.GetUsagePercentage(TaskMessage)
	searchPercent, _ := rl.GetUsagePercentage(TaskSearch)

	stats := fmt.Sprintf(`Daily Rate Limit Usage:
  Connections: %d/%d (%.1f%%)
  Messages:    %d/%d (%.1f%%)
  Searches:    %d/%d (%.1f%%)
  Resets at:   %s`,
		limit.ConnectionCount, rl.config.MaxConnectionsPerDay, connPercent,
		limit.MessageCount, rl.config.MaxMessagesPerDay, msgPercent,
		limit.SearchCount, rl.config.MaxSearchesPerDay, searchPercent,
		rl.getNextMidnight().Format("15:04:05"))

	return stats, nil
}

// CanPerformTask checks if a task can be performed (combines limit check and cooldown)
func (rl *RateLimiter) CanPerformTask(taskType TaskType) error {
	// Check daily limit
	if err := rl.CheckDailyLimit(taskType); err != nil {
		return err
	}

	// Warn if approaching limit
	shouldWarn, _ := rl.ShouldWarnAboutLimit(taskType)
	if shouldWarn {
		remaining, _ := rl.GetRemainingQuota(taskType)
		logger.Warning(fmt.Sprintf("Approaching rate limit for %s: %d actions remaining", taskType, remaining))
	}

	return nil
}

package stealth

import (
	"testing"
	"time"
)

// TestRandomDelay verifies that delays are within expected range
func TestRandomDelay(t *testing.T) {
	start := time.Now()
	RandomDelay(100, 200)
	elapsed := time.Since(start)

	// Should take between 100-200ms
	if elapsed < 100*time.Millisecond {
		t.Errorf("Delay too short: %v", elapsed)
	}
	if elapsed > 250*time.Millisecond { // Allow some margin
		t.Errorf("Delay too long: %v", elapsed)
	}
}

// TestRandomDelayMultipleCalls ensures randomness
func TestRandomDelayMultipleCalls(t *testing.T) {
	delays := make([]time.Duration, 5)

	for i := 0; i < 5; i++ {
		start := time.Now()
		RandomDelay(50, 150)
		delays[i] = time.Since(start)
	}

	// At least some variation in delays
	minDelay := delays[0]
	maxDelay := delays[0]

	for _, d := range delays {
		if d < minDelay {
			minDelay = d
		}
		if d > maxDelay {
			maxDelay = d
		}
	}

	if minDelay == maxDelay {
		t.Error("Delays should vary - possible issue with randomization")
	}
}

// TestRandomDelayWithSmallRange tests with very small delay range
func TestRandomDelayWithSmallRange(t *testing.T) {
	start := time.Now()
	RandomDelay(10, 20)
	elapsed := time.Since(start)

	if elapsed < 10*time.Millisecond {
		t.Errorf("Delay too short: %v", elapsed)
	}
	if elapsed > 30*time.Millisecond {
		t.Errorf("Delay too long: %v", elapsed)
	}
}

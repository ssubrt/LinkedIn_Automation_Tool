package stealth

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// MoveMouseRandomly simulates small human-like mouse movements to avoid detection.
// It performs multiple random mouse movements across the page with natural pauses
// to mimic real human behavior patterns.
func MoveMouseRandomly(page *rod.Page) {
	// Create a seeded random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Perform 3-5 random mouse movements
	numMovements := 3 + r.Intn(3) // Random number between 3-5

	for i := 0; i < numMovements; i++ {
		// Generate random X coordinate between 100-800 pixels (wider range for natural movement)
		x := r.Intn(700) + 100
		// Generate random Y coordinate between 100-600 pixels
		y := r.Intn(500) + 100

		// Move the mouse to the generated coordinates
		page.Mouse.MustMoveTo(float64(x), float64(y))
		// Pause between 300-800ms for each movement
		time.Sleep(time.Duration(300+r.Intn(500)) * time.Millisecond)
	}
}

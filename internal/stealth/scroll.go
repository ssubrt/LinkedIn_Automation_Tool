package stealth

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// RandomScroll simulates human-like scrolling behavior on a webpage.
// It performs multiple scrolls with random distances and pauses to mimic natural browsing patterns.
func RandomScroll(page *rod.Page) {
	// Create a seeded random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Perform 3-5 random scroll actions
	numScrolls := 3 + r.Intn(3) // Random number between 3-5

	for i := 0; i < numScrolls; i++ {
		// Generate a random scroll distance between 200-600 pixels vertically
		scrollY := r.Intn(400) + 200

		// Scroll the page vertically (0 = no horizontal scroll, scrollY = vertical scroll distance)
		page.Mouse.MustScroll(0, float64(scrollY))
		// Pause for 800-1500ms to simulate human reading time between scrolls
		time.Sleep(time.Duration(800+r.Intn(700)) * time.Millisecond)
	}
}

package stealth

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// Point represents a 2D coordinate
type Point struct {
	X float64
	Y float64
}

// MoveBezier moves the mouse along a Bézier curve from start to end point
// This creates natural, human-like mouse movements instead of straight lines
func MoveBezier(page *rod.Page, fromX, fromY, toX, toY float64) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate random control points for the Bézier curve
	// Control points determine the curve's shape
	cp1X := fromX + (toX-fromX)*0.25 + float64(r.Intn(100)-50)
	cp1Y := fromY + (toY-fromY)*r.Float64()

	cp2X := fromX + (toX-fromX)*0.75 + float64(r.Intn(100)-50)
	cp2Y := fromY + (toY-fromY)*r.Float64()

	// Number of steps in the curve (20-30 for smooth movement)
	steps := 20 + r.Intn(11)

	// Move mouse along the Bézier curve
	for i := 0; i <= steps; i++ {
		// Calculate parameter t (0 to 1)
		t := float64(i) / float64(steps)

		// Apply easing for variable speed (slow start/end, fast middle)
		t = easeInOutCubic(t)

		// Cubic Bézier curve formula: B(t) = (1-t)³P₀ + 3(1-t)²tP₁ + 3(1-t)t²P₂ + t³P₃
		x := math.Pow(1-t, 3)*fromX +
			3*math.Pow(1-t, 2)*t*cp1X +
			3*(1-t)*math.Pow(t, 2)*cp2X +
			math.Pow(t, 3)*toX

		y := math.Pow(1-t, 3)*fromY +
			3*math.Pow(1-t, 2)*t*cp1Y +
			3*(1-t)*math.Pow(t, 2)*cp2Y +
			math.Pow(t, 3)*toY

		// Move to calculated position
		page.Mouse.MustMoveTo(x, y)

		// Small delay between movements (1-3ms for smooth animation)
		time.Sleep(time.Duration(1+r.Intn(3)) * time.Millisecond)
	}

	// Add slight overshoot and correction (human behavior)
	if r.Float64() > 0.7 {
		overshootX := toX + float64(r.Intn(10)-5)
		overshootY := toY + float64(r.Intn(10)-5)
		page.Mouse.MustMoveTo(overshootX, overshootY)
		time.Sleep(time.Duration(10+r.Intn(20)) * time.Millisecond)
		page.Mouse.MustMoveTo(toX, toY)
	}
}

// easeInOutCubic provides natural acceleration/deceleration
// Slow at start and end, fast in the middle (like human movement)
func easeInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}

// MoveMouseRandomly simulates small human-like mouse movements to avoid detection.
// It performs multiple random mouse movements across the page with natural pauses
// to mimic real human behavior patterns.
func MoveMouseRandomly(page *rod.Page) {
	// Create a seeded random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Get current mouse position (or start from a random position)
	currentX := float64(200 + r.Intn(400))
	currentY := float64(150 + r.Intn(300))

	// Perform 3-5 random mouse movements
	numMovements := 3 + r.Intn(3) // Random number between 3-5

	for i := 0; i < numMovements; i++ {
		// Generate random target coordinates
		targetX := float64(r.Intn(700) + 100) // 100-800 pixels
		targetY := float64(r.Intn(500) + 100) // 100-600 pixels

		// Move using Bézier curve for natural movement
		MoveBezier(page, currentX, currentY, targetX, targetY)

		// Update current position
		currentX = targetX
		currentY = targetY

		// Pause between movements (300-800ms)
		time.Sleep(time.Duration(300+r.Intn(500)) * time.Millisecond)
	}
}

// HoverRandomElements hovers the mouse over random interactive elements on the page
// This simulates natural browsing behavior where users hover over links and buttons
func HoverRandomElements(page *rod.Page) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Find all interactive elements (links, buttons)
	elements, err := page.Elements("a, button, [role='button']")
	if err != nil || len(elements) == 0 {
		// If no elements found, just do random movements
		MoveMouseRandomly(page)
		return nil
	}

	// Hover over 2-3 random elements
	numHovers := 2 + r.Intn(2)
	if numHovers > len(elements) {
		numHovers = len(elements)
	}

	// Shuffle and select random elements
	shuffled := make([]*rod.Element, len(elements))
	copy(shuffled, elements)
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	for i := 0; i < numHovers; i++ {
		element := shuffled[i]

		// Get element position
		shape, err := element.Shape()
		if err != nil {
			continue
		}

		// Get first quad (box) from shape
		if len(shape.Quads) == 0 {
			continue
		}
		quad := shape.Quads[0]

		// Calculate center of element
		centerX := (quad[0] + quad[2] + quad[4] + quad[6]) / 4
		centerY := (quad[1] + quad[3] + quad[5] + quad[7]) / 4

		// Get current mouse position (approximate)
		currentX := float64(200 + r.Intn(400))
		currentY := float64(150 + r.Intn(300))

		// Move to element with Bézier curve
		MoveBezier(page, currentX, currentY, centerX, centerY)

		// Hover for 200-500ms (simulating user reading/thinking)
		time.Sleep(time.Duration(200+r.Intn(300)) * time.Millisecond)
	}

	return nil
}

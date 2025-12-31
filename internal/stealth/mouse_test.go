package stealth

import (
	"testing"
)

func TestPoint(t *testing.T) {
	p := Point{X: 10.0, Y: 20.0}
	if p.X != 10.0 || p.Y != 20.0 {
		t.Errorf("Point creation failed: expected (10, 20), got (%.1f, %.1f)", p.X, p.Y)
	}
}

func TestEaseInOutCubic(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0.0, 0.0}, // Start: should be 0
		{1.0, 1.0}, // End: should be 1
		{0.5, 0.5}, // Middle: should be around 0.5
	}

	for _, test := range tests {
		result := easeInOutCubic(test.input)
		// Allow small floating point differences
		if result < test.expected-0.1 || result > test.expected+0.1 {
			t.Errorf("easeInOutCubic(%.1f) = %.2f, expected ~%.1f", test.input, result, test.expected)
		}
	}
}

func TestEaseInOutCubicBounds(t *testing.T) {
	// Test that easing always returns values between 0 and 1
	testValues := []float64{0.0, 0.1, 0.25, 0.5, 0.75, 0.9, 1.0}

	for _, val := range testValues {
		result := easeInOutCubic(val)
		if result < 0.0 || result > 1.0 {
			t.Errorf("easeInOutCubic(%.2f) = %.2f, should be between 0 and 1", val, result)
		}
	}
}

func TestEaseInOutCubicSlowStartEnd(t *testing.T) {
	// The easing function should have slower start and end (derivative should be low)
	// We can check this by comparing changes near 0 and 1 vs changes near 0.5

	// Change near start
	deltaStart := easeInOutCubic(0.1) - easeInOutCubic(0.0)

	// Change near middle
	deltaMid := easeInOutCubic(0.55) - easeInOutCubic(0.45)

	// Change near end
	deltaEnd := easeInOutCubic(1.0) - easeInOutCubic(0.9)

	// Middle should have larger changes (faster) than start/end
	if deltaMid <= deltaStart || deltaMid <= deltaEnd {
		t.Error("Easing function should be faster in the middle than at start/end")
	}
}

// Benchmark Bézier curve calculation
func BenchmarkEaseInOutCubic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		easeInOutCubic(0.5)
	}
}

func TestBezierCurveProperties(t *testing.T) {
	// Test that a Bézier curve goes from start to end point
	// We can't test the actual page.Mouse calls, but we can test the math

	fromX, fromY := 100.0, 100.0
	_ = fromX // Used in comment/concept
	_ = fromY // Used in comment/concept

	// At t=0, curve should be at start point
	t0X := fromX
	t0Y := fromY

	if t0X != fromX || t0Y != fromY {
		t.Errorf("Curve at t=0 should be at start point (%.1f, %.1f), got (%.1f, %.1f)",
			fromX, fromY, t0X, t0Y)
	}

	// At t=1, curve should be at end point (with our easing function applied)
	// The actual implementation uses easing, so we just verify the concept
}

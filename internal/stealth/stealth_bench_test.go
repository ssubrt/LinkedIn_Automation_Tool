package stealth

import (
	"math/rand"
	"testing"
	"time"
)

// BenchmarkRandomDelay measures delay performance
func BenchmarkRandomDelay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomDelay(10, 20)
	}
}

// BenchmarkRandomNumber measures random number generation performance
func BenchmarkRandomNumber(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Intn(1000)
	}
}

// BenchmarkMoveMouseRandomly measures mouse movement performance
func BenchmarkMoveMouseRandomly(b *testing.B) {
	// Note: This benchmark requires a valid page object
	// In real scenario, you would mock the page
	// For now, we just measure the randomization logic
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		numMovements := 3 + r.Intn(3)
		for j := 0; j < numMovements; j++ {
			_ = r.Intn(700) + 100
			_ = r.Intn(500) + 100
		}
	}
}

// BenchmarkRandomScroll measures scroll generation performance
func BenchmarkRandomScroll(b *testing.B) {
	// Benchmark scroll distance calculation
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		numScrolls := 3 + r.Intn(3)
		for j := 0; j < numScrolls; j++ {
			_ = r.Intn(400) + 200
		}
	}
}

// BenchmarkRandomNumberRange tests various range sizes
func BenchmarkRandomNumberRange100(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Intn(100)
	}
}

func BenchmarkRandomNumberRange1000(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Intn(1000)
	}
}

func BenchmarkRandomNumberRange10000(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Intn(10000)
	}
}

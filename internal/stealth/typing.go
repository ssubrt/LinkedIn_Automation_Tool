package stealth

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// TypeLikeHuman types text character by character with random delays
func TypeLikeHuman(el *rod.Element, text string) {
	for _, char := range text {
		el.MustInput(string(char))

		time.Sleep(time.Duration(100+rand.Intn(150)) * time.Millisecond)
	}
}

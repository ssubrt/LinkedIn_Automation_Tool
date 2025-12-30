package stealth

// math/rand used to generate random delays / nummbers
//time used to pause executions  for a certain duration
import (
	"math/rand"
	"time"
) // both these are required to behave like human

// minMs and maxMs : delay in Milliseconds
// eg: RandomDelay(1000,500)  will generate a random delay between 500ms to 1000ms
func RandomDelay(minMs int, maxMs int) {

	//seed initializes the random generator
	//time.Now().UnixNano() provides a different seed value each time the program runs
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	delay := r.Intn(maxMs-minMs+1) + minMs

	time.Sleep(time.Duration(delay) * time.Millisecond)

}

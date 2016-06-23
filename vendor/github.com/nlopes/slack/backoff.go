package slack

import (
	"math"
	"math/rand"
	"time"
)

// This one was ripped from https://github.com/jpillora/backoff/blob/master/backoff.go

// Backoff is a time.Duration counter. It starts at Min.  After every
// call to Duration() it is multiplied by Factor.  It is capped at
// Max. It returns to Min on every call to Reset().  Used in
// conjunction with the time package.
type backoff struct {
	attempts int
	//Factor is the multiplying factor for each increment step
	Factor float64
	//Jitter eases contention by randomizing backoff steps
	Jitter bool
	//Min and Max are the minimum and maximum values of the counter
	Min, Max time.Duration
}

// Returns the current value of the counter and then multiplies it
// Factor
func (b *backoff) Duration() time.Duration {
	//Zero-values are nonsensical, so we use
	//them to apply defaults
	if b.Min == 0 {
		b.Min = 100 * time.Millisecond
	}
	if b.Max == 0 {
		b.Max = 10 * time.Second
	}
	if b.Factor == 0 {
		b.Factor = 2
	}
	//calculate this duration
	dur := float64(b.Min) * math.Pow(b.Factor, float64(b.attempts))
	if b.Jitter == true {
		dur = rand.Float64()*(dur-float64(b.Min)) + float64(b.Min)
	}
	//cap!
	if dur > float64(b.Max) {
		return b.Max
	}
	//bump attempts count
	b.attempts++
	//return as a time.Duration
	return time.Duration(dur)
}

//Resets the current value of the counter back to Min
func (b *backoff) Reset() {
	b.attempts = 0
}

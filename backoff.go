package throttle

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

const (
	// DefaultMaxWait default wait time for backoff.
	DefaultMaxWait = 15 * time.Minute
)

// Backoff is an exponential backoff wrapper that runs a function at an random and exponentially
// growing time until the function is successful or the MaxWait time is hit.
type Backoff struct {
	MaxWait      time.Duration
	BackoffCodes []int
	random       *rand.Rand
}

// ExpBackoff creates and instance of Backoff using default values.
func ExpBackoff() *Backoff {
	return &Backoff{
		MaxWait:      DefaultMaxWait,
		BackoffCodes: []int{},
		random:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SetMaxWait adds duration to MaxWait in *Backoff.
func (bo *Backoff) SetMaxWait(dur time.Duration) {
	bo.MaxWait = dur
}

// Run runs function until return nil error or MaxWait is hit.
func (bo *Backoff) Run(f func() error) error {
	var wait time.Duration
	attempts := 0.0
	for wait < bo.MaxWait {
		err := f()
		if err == nil {
			return nil
		}
		if noGo, ok := err.(*NoGosErrors); ok {
			return &Error{`Run`, `hit nogo error`, noGo.Err}
		}
		wait = time.Duration(int(math.Pow(2, attempts)+float64(bo.random.Intn(1000)))) * time.Millisecond
		time.Sleep(wait)
		attempts++
	}
	return &Error{`Run`, `tried function`, errors.New("max wait time hit")}
}

// NoGosErrors struct holds errors that will fail backoff.
// Found this way of doing errors from cenkalti on github,
// they have an exponential backoff package that helped me
// create this one. You can find it at  github.com/cenkalti/backoff.
type NoGosErrors struct {
	Err error
}

func (ng *NoGosErrors) Error() string {
	return ng.Err.Error()
}

// NoGos warps err in *NoGos.
// Use in function passed into Run()
func NoGos(err error) *NoGosErrors {
	return &NoGosErrors{err}
}

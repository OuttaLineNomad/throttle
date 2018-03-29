package throttle

import (
	"errors"
	"time"
)

// Bucket for LeakyBucket() set Size, Rate, and channel to control flow
type Bucket struct {
	Size      int
	Rate      time.Duration
	Regulator chan bool
}

// LeakyBucket creates a Bucket with a rate to regulate api calls.
func LeakyBucket(size int, rate time.Duration) *Bucket {
	buk := &Bucket{
		Size:      size,
		Regulator: make(chan bool, size),
		Rate:      rate,
	}

	for i := 0; i < cap(buk.Regulator); i++ {
		buk.Regulator <- true
	}

	go (func(buk *Bucket, rate time.Duration) {
		drop := time.Tick(rate)
		for {
			bukLen := len(buk.Regulator)
			select {
			case <-drop:
				if bukLen < cap(buk.Regulator) {
					buk.Regulator <- true
				}
			}
		}
	})(buk, buk.Rate)

	return buk
}

// Take takes one token from Bucket if there is one to take, if not it waits until it can take one.
func (buk *Bucket) Take() {
	<-buk.Regulator
}

// TakeN takes n amount of tokens, if n amount is not found thant waits until it can take all.
// n cannot exceed the amount of Bucket size.
func (buk *Bucket) TakeN(n int) error {
	if n > buk.Size {
		return &Error{"TakeN", "taking n tokens out of bucket", errors.New("n exceeds the amount of bucket size")}
	}
	lenReg := len(buk.Regulator)
	if n <= lenReg {
		for i := 0; i < n; i++ {
			<-buk.Regulator
		}
		return nil
	}

	needed := n - lenReg
	wait := time.Duration(needed) * buk.Rate
	time.Sleep(wait)
	for i := 0; i < n; i++ {
		<-buk.Regulator
	}
	return nil
}

// GetRate creates time duration for calls.
// Determines how many calls per time given.
func GetRate(calls int, dur time.Duration) time.Duration {
	return dur / time.Duration(calls)
}
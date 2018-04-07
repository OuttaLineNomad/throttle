// Copyright 2018 Bryce Mullen. All rights reserved.
// Use of this source code is governed by a Apache-2.0
// license that can be found in the LICENSE file.
//
package throttle_test

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/OuttaLineNomad/throttle"
)

// LeakyBucket sets size 5 bucket. Meaning only 5 calls can be made in a burst.
// After that the bucket is regulated at the rate of of 15 calls a minute. GetRate
// does the math for you.
//
// README:Leaky Bucket
func ExampleLeakyBucket() {
	// GetRate is used to get the duration to regulate the bucket.
	bkt := throttle.LeakyBucket(5, throttle.GetRate(15, time.Minute))
	for i := 0; i < 6; i++ {
		bkt.Take()
		// after this you call you make your call.
		resp, err := http.Get("http://httpbin.org/get")
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(resp.StatusCode)
	}
	// Output:
	// 200
	// 200
	// 200
	// 200
	// 200
	// 200
}

// README:Leaky Bucket
func ExampleLeakyBucket_second() {
	// GetRate is used to get the duration to regulate the bucket.
	bkt := throttle.LeakyBucket(5, throttle.GetRate(15, time.Minute))
	for i := 0; i < 6; i++ {
		bkt.Take()
		// after this you call you make your call.
		resp, err := http.Get("http://httpbin.org/get")
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(resp.StatusCode)
	}
	// Output:
	// 200
	// 200
	// 200
	// 200
	// 200
	// 200
}

// README:Leaky Bucket
func ExampleLeakyBucket_third() {
	atOnce := 5
	// GetRate is used to get the duration to regulate the bucket.
	bkt := throttle.LeakyBucket(5, throttle.GetRate(15, time.Minute))
	for i := 0; i < 2; i++ {
		bkt.TakeN(atOnce)
		for j := 0; j < atOnce; j++ {
			// after this you call you make your call.
			resp, err := http.Get("http://httpbin.org/get")
			if err != nil {
				log.Panic(err)
			}
			fmt.Println(resp.StatusCode)
		}
	}
	// Output:
	// 200
	// 200
	// 200
	// 200
	// 200
	// 200
	// 200
	// 200
	// 200
	// 200
}

// ExpBackoff is an exponential backoff that runs a function at a random and exponentially
// growing time until the function is successful or the MaxWait time is hit.
//
// README:ExpBackoff
func ExampleExpBackoff() {
	backOff := throttle.ExpBackoff()
	var myResp *http.Response
	apiFunc := func() error {
		resp, err := http.Get("http://httpbin.org/get")
		myResp = resp
		return err
	}
	backOff.Run(apiFunc)
	fmt.Println(myResp.StatusCode)
	// Output: 200
}

// This example uses SetMaxWait function which sets the max wait duration. Default is 15 minutes.
// Also it uses NoGos. Which is for errors from the API that you don't want your application to ignore.
// By default function retry when receives any error.
//
// README:Advanced ExpBackoff
func ExampleExpBackoff_second() {
	backOff := throttle.ExpBackoff()
	backOff.SetMaxWait(5 * time.Minute)
	var myResp *http.Response
	apiFunc := func() error {
		resp, err := http.Get("http://httpbin.org/get")
		if err != nil {
			if err.Error() == "Error I don't want" {
				log.Println("This is a NoGo error.")
				return throttle.NoGos(err)
			}
			return err
		}
		myResp = resp
		return nil
	}
	backOff.Run(apiFunc)
	fmt.Println(myResp.StatusCode)
	// Output:200
}

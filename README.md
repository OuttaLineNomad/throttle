# Throttle 

## Overview

Package throttle allows users to use functions to keep with different API throttling policies.

Use LeakyBucket to make calls to web services that require a leaky bucket algorithm, and use ExpBackoff if you need an exponential backoff call.

## Install

```
go get github.com/OuttaLineNomad/throttle 
```

## Examples
### Using LeakyBucket
LeakyBucket sets size 4 bucket. Meaning only 10 calls can be made in a burst. After that the bucket is regulated at the rate of of 4 calls a minute. GetRate does the math for you.
```go
func main() {
	bkt := throttle.LeakyBucket(10, throttle.GetRate(4, time.Minute))
	for i := 0; i < 100; i++ {
		bkt.Take()
		// after this you call you make your call.
		resp, err := http.Get("http://httpbin.org/get")
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(resp.StatusCode)
	}
}
```
### Using ExpBackoff
ExpBackoff is an exponential backoff that runs a function at a random and exponentially growing time until the function is successful or the MaxWait time is hit.

#### Simple Example
```go
func main() {
	backOff := throttle.ExpBackoff()
	var myResp *http.Response
	apiFunc := func() error {
		resp, err := http.Get("http://httpbin.org/get")
		myResp = resp
		return err
	}
	backOff.Run(apiFunc)
	fmt.Println(myResp.StatusCode)
}
```

#### Advanced Example

This example uses SetMaxWait function which sets the max wait duration. Default is 15 minutes. Also it uses NoGos. Which is for errors from the API that you don't want your application to ignore. By default function retry when receives any error.
```go
func main() {
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
}

```

## Author

* **Bryce Mullen** - *Project Manager*  @Wedgenix connect [here](https://www.linkedin.com/in/bryce-mullen).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

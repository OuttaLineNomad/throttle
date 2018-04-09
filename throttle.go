// Package throttle helps control API calls to webservices.
// Allowing programmers to choose different ways to throttle their calls
// to adhere to the throttling policies of certain APIs.
//
// API calls to a webservice are generally throttled. Different policies
// are used, this package helps follow the most used polices.
//
// Leaky Bucket
//
// Leaky Bucket or Token Bucket Algorithm is when a service sets a limit
// on a steady rate and bust request amount. The burst is the size of the
// bucket, and the rate is how cast that bucket will fill. When API calls
// exceed the steady rate and bucket size the service sends some type of
// error.
//
// LeakyBucket function takes two parameters:
//      size: size of bucket of burst requests allowed
//      rate: what rate the bucket fills.
//
// Exponential Backoff
//
// This is less of a webservice throttling policy and more of calling policy.
// Exponential Backoff calls API at a random and exponentially growing time until
// the call is successful or the MaxWait time is hit.
//
// ExpBackoff is a method. Start it then use Run with one parameter:
//      f: f is a funciton that returns an error. See examples for more details.
// 	MaxWait: the max time the function will keep retrying the call.
//
package throttle

// Error type to share error messages from package.
type Error struct {
	Msg  string
	Func string
	Err  error
}

func (er *Error) Error() string {
	return `throttle.` + er.Func + `: ` + er.Msg + ` : ` + er.Err.Error()
}

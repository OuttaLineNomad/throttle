// Package throttle allows users to use different functions to
// keep with different api throttling policies.
//
// Use LeakyBucket to make calls to web services that
// require a leaky bucket algorithm.
//
// Use ExpBackoff if you need an exponential backoff call.
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

package proxy

import "net/http"

// LBOption is a container for optional properties that can be used for initializing the Balancer.
type LBOption struct {
	selector Selector
}

// WithSelector sets the Selector to use for the Balancer.
func WithSelector(selector Selector) func(*LBOption) {
	return func(o *LBOption) {
		o.selector = selector
	}
}

// HostOption is a container for optional properties that can be used for initializing a Host.
type HostOption struct {
	errorHandler func(http.ResponseWriter, *http.Request, error)
	transport    http.RoundTripper
}

// WithTransport sets the http.RoundTripper transport for a Host.
func WithTransport(transport http.RoundTripper) func(*HostOption) {
	return func(o *HostOption) {
		o.transport = transport
	}
}

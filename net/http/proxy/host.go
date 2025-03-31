package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/transientvariable/anchor"
)

const (
	StatusClientClosedRequest     = 499
	StatusClientClosedRequestText = "Client Closed Request"
)

// Host defines the attributes and behavior for a network proxy host.
type Host struct {
	failures      int
	inactive      bool
	inactiveSince time.Time
	proxy         *httputil.ReverseProxy
	mutex         sync.RWMutex
	target        *url.URL
}

// NewHost creates a new Host from the provided address string and options.
func NewHost(target string, options ...func(*HostOption)) (*Host, error) {
	t, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("proxy_host: failed to parse target URL %v: %w", t, err)
	}
	h := &Host{target: t, proxy: httputil.NewSingleHostReverseProxy(t)}

	opts := &HostOption{}
	for _, opt := range options {
		opt(opts)
	}
	return h, nil
}

// Active returns whether the Host is active.
func (h *Host) Active() bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return !h.inactive
}

// Failures returns the number of failures for the Host.
func (h *Host) Failures() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	f := h.failures
	return f
}

// InactiveSince returns the timestamp indicating the last time the Host was active.
func (h *Host) InactiveSince() time.Time {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	s := h.inactiveSince
	return s
}

// Target returns the Host url.URL target.
func (h *Host) Target() (*url.URL, error) {
	if h.target == nil {
		return nil, errors.New("proxy_host: address not set")
	}
	addr, err := url.Parse(h.target.String())
	if err != nil {
		return addr, fmt.Errorf("proxy_host: invalid address %q: %v", h.target.String(), err)
	}
	return addr, nil
}

// String returns a string representation of the Host attributes.
func (h *Host) String() string {
	return string(anchor.ToJSON(h.toMap()))
}

// markActive sets the Host active status to true
func (h *Host) markActive() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.inactive = false
}

// markHealthy marks the Host as healthy.
func (h *Host) markHealthy() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.inactive = false
	h.inactiveSince = time.Time{}
	h.failures = 0
}

// markInactive marks the Host as inactive.
func (h *Host) markInactive() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.inactive = true
	if h.inactiveSince.IsZero() {
		h.inactiveSince = time.Now().UTC()
	}
	h.failures++
}

// serveHTTP performs the request for the Host.
func (h *Host) serveHTTP(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}

// toMap returns a map representing the Host attributes.
func (h *Host) toMap() map[string]any {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	m := make(map[string]any)
	if h.target != nil {
		m["target"] = h.target.String()
	}
	m["active"] = h.Active()
	m["failures"] = h.failures
	return m
}

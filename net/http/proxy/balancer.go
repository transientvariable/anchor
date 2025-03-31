package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/transientvariable/anchor"
	"github.com/transientvariable/log-go"
)

// Selector defines the behavior for selecting a Host proxy from a Pool.
type Selector interface {
	Select(...*Host) (*Host, error)
}

type roundRobinSelector struct {
	current int
	mutex   sync.RWMutex
}

// Select returns a Host proxy in a round-robin manner.
func (s *roundRobinSelector) Select(hosts ...*Host) (*Host, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.current = (s.current + 1) % len(hosts)

	h := hosts[s.current]

	log.Trace("[proxy:balancer] selected host", log.Int("index", s.current))

	return h, nil
}

// Balancer defines the behavior for multiplexing HTTP requests amongst of a number of Host proxies.
type Balancer interface {
	http.Handler

	// Targets returns the list of URLs of available Host proxies.
	Targets() ([]*url.URL, error)
}

type balancer struct {
	active                 []*Host
	inactive               []*Host
	mutex                  sync.RWMutex
	reviveTimeout          time.Duration
	reviveTimeoutThreshold int
	selector               Selector
}

// NewBalancer creates a new proxy Balancer using the provided Selector and Host proxy list.
//
// If the provided Selector is nil, a default one based the round-robin algorithm is used.
func NewBalancer(hosts []*Host, options ...func(*LBOption)) (Balancer, error) {
	l := &balancer{
		active:   []*Host{},
		inactive: []*Host{},
	}

	// sanitize the list of provided hosts and add them to the load balancer as active hosts
	for _, h := range hosts {
		if h != nil {
			t, err := h.Target()
			if err != nil {
				return nil, fmt.Errorf("load_balancer: %w", err)
			}
			log.Debug("[proxy:balancer] adding host", log.String("target", t.String()))
			l.active = append(l.active, h)
		}
	}

	// we need at least one host to load balance
	if len(l.active) == 0 {
		return nil, errors.New("load_balancer: at least one host must be provided")
	}

	opts := &LBOption{}
	for _, opt := range options {
		opt(opts)
	}

	if opts.selector == nil {
		// if the option to set a Selector is nil, use the round-robin selector as the default
		l.selector = &roundRobinSelector{current: -1}
	}
	log.Debug(fmt.Sprintf("[proxy:balancer]: \n%s", l))
	return l, nil
}

// ServeHTTP performs the HTTP request using one of the active Host proxies of the Balancer.
func (b *balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.active) > 0 {
		h, err := b.selector.Select(b.active...)
		if err != nil {
			http.Error(w, "Service not available", http.StatusServiceUnavailable)
			return
		}
		h.serveHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

// Targets returns the list of URLs of available Host proxies for the Balancer.
func (b *balancer) Targets() ([]*url.URL, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	var t []*url.URL
	for _, h := range b.active {
		a, err := h.Target()
		if err != nil {
			return t, err
		}
		t = append(t, a)
	}
	return t, nil
}

// String returns a string representation of the Balancer.
func (b *balancer) String() string {
	return string(anchor.ToJSON(b.toMap()))
}

// toMap returns a map representing the balancer attributes.
func (b *balancer) toMap() map[string]any {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	m := make(map[string]any)

	var activeHosts []map[string]any
	for _, h := range b.active {
		activeHosts = append(activeHosts, h.toMap())
	}

	var inactiveHosts []map[string]any
	for _, h := range b.inactive {
		inactiveHosts = append(inactiveHosts, h.toMap())
	}

	m["pool"] = map[string]any{
		"hosts": len(b.active) + len(b.inactive),
		"active": map[string]any{
			"count": len(activeHosts),
			"hosts": activeHosts,
		},
		"inactive": map[string]any{
			"count": len(inactiveHosts),
			"hosts": inactiveHosts,
		},
	}
	return m
}

func writeStatus(w http.ResponseWriter, sc int) (string, error) {
	w.WriteHeader(sc)
	var st string
	if sc == http.StatusMisdirectedRequest || sc == StatusClientClosedRequest {
		st = StatusClientClosedRequestText
	} else {
		st = http.StatusText(sc)
	}

	_, err := w.Write([]byte(st))
	if err != nil {
		return st, err
	}
	return st, nil
}

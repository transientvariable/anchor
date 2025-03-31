package http

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"runtime"
	"time"

	"github.com/transientvariable/anchor"
	"github.com/transientvariable/log-go"

	"github.com/cenkalti/backoff/v4"

	gohttp "net/http"
)

const (
	BufferSizeRead    = 4 * anchor.KiB
	BufferSizeWrite   = 4 * anchor.KiB
	ConnIdleTimeout   = 90 * time.Second
	ConnMaxPerHost    = 512
	ConnMaxIdle       = 100
	DialKeepAlive     = 30 * time.Second
	DialTimeout       = 30 * time.Second
	DisableKeepAlives = true

	// RetryMax sets the number times an HTTP request is retried.
	RetryMax = 2

	Timeout = 90 * time.Second
)

var (
	errRedirectsPattern  = regexp.MustCompile(`stopped after \d+ redirects\z`)
	errSchemePattern     = regexp.MustCompile(`unsupported protocol scheme`)
	errNotTrustedPattern = regexp.MustCompile(`certificate is not trusted`)
)

// NewClient ...
func NewClient() *gohttp.Client {
	return &gohttp.Client{
		Timeout:   Timeout,
		Transport: NewTransport(),
	}
}

// DefaultClient returns a new http.Client with similar default values to http.Client, but with a non-shared
// http.Transport which has keep-alives disabled.
func DefaultClient() *gohttp.Client {
	return &gohttp.Client{
		Timeout:   Timeout,
		Transport: DefaultTransport(),
	}
}

// NewTransport ...
func NewTransport() *gohttp.Transport {
	transport := DefaultTransport()
	transport.DialContext = (&net.Dialer{
		KeepAlive: DialKeepAlive,
		Timeout:   DialTimeout,
	}).DialContext
	transport.IdleConnTimeout = ConnIdleTimeout
	transport.MaxConnsPerHost = ConnMaxPerHost
	transport.MaxIdleConns = ConnMaxIdle
	transport.ReadBufferSize = BufferSizeRead
	transport.WriteBufferSize = BufferSizeWrite
	transport.DisableKeepAlives = DisableKeepAlives
	return transport
}

// DefaultTransport ...
func DefaultTransport() *gohttp.Transport {
	return &gohttp.Transport{
		DialContext: (&net.Dialer{
			Timeout:   DialTimeout,
			KeepAlive: DialKeepAlive,
		}).DialContext,
		ExpectContinueTimeout: 3 * time.Second,
		ForceAttemptHTTP2:     true,
		IdleConnTimeout:       ConnIdleTimeout,
		MaxConnsPerHost:       ConnMaxPerHost,
		MaxIdleConns:          ConnMaxIdle,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		Proxy:                 gohttp.ProxyFromEnvironment,
		ReadBufferSize:        BufferSizeRead,
		TLSHandshakeTimeout:   10 * time.Second,
		WriteBufferSize:       BufferSizeWrite,
	}
}

// DoWithRetry ...
func DoWithRetry(client *gohttp.Client, req *gohttp.Request) (*gohttp.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	if req.Context().Err() != nil {
		return resp, req.Context().Err()
	}

	if retry, _ := Retry(resp, err); retry {
		err = backoff.Retry(func() error {
			log.Trace("[anchor] retry", log.String("url", req.URL.String()))
			rr, err := client.Do(req)
			if err != nil {
				if rr.Body != nil {
					if err := rr.Body.Close(); err != nil {
						log.Error("[anchor] retry", log.Err(err))
					}
				}
				return err
			}
			resp = rr
			return nil
		}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), RetryMax))
	}
	return resp, err
}

// Retry ...
func Retry(resp *gohttp.Response, err error) (bool, error) {
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			if errors.Is(urlErr.Err, x509.UnknownAuthorityError{}) {
				return false, err
			}

			e := urlErr.Error()
			if errRedirectsPattern.MatchString(e) || errSchemePattern.MatchString(e) || errNotTrustedPattern.MatchString(e) {
				return false, urlErr
			}
		}
		return true, nil
	}

	if resp.StatusCode == gohttp.StatusTooManyRequests {
		return true, nil
	}

	if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != gohttp.StatusNotImplemented) {
		return true, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}
	return false, nil
}

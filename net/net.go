package net

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// NewSOCKS5DialContext creates a new proxy.ContextDialer for SOCKS5 proxy.
func NewSOCKS5DialContext(proxyURL *url.URL) (func(context.Context, string, string) (net.Conn, error), error) {
	if proxyURL == nil {
		return nil, fmt.Errorf("socks5: proxy addr is required")
	}

	conn, err := proxy.SOCKS5("tcp", proxyURL.String(), nil, &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("socks5: %w", err)
	}

	if ctxDialer, ok := conn.(proxy.ContextDialer); ok {
		return ctxDialer.DialContext, nil
	}
	return nil, errors.New("socks5: could not assert context dialer from proxy")
}

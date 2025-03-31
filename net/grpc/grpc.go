package grpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/transientvariable/anchor/net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	gonet "net"
)

const (
	MinKeepAliveTime    = 15 * time.Second
	MinKeepAliveTimeout = 30 * time.Second
)

// New creates a new GRPC connection.
func New(target string, options ...func(*Option)) (*grpc.ClientConn, error) {
	if target = strings.TrimSpace(target); target == "" {
		return nil, fmt.Errorf("grpc: target is required")
	}

	opts := &Option{}
	for _, opt := range options {
		opt(opts)
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(int(opts.messageSizeMaxReceive)),
			grpc.MaxCallSendMsgSize(int(opts.messageSizeMaxSend)),
		),
	}

	if opts.socks5Enabled {
		var address string
		if address = strings.TrimSpace(os.Getenv("HTTP_PROXY")); address == "" {
			address = strings.TrimSpace(os.Getenv("HTTPS_PROXY"))
			return nil, fmt.Errorf("grpc: SOCKS5 enabled but proxy has not been set")
		}

		proxyURL, err := url.Parse(address)
		if err != nil {
			return nil, fmt.Errorf("grpc: could not parse proxy addr: %s: %w", address, err)
		}

		dialCtx, err := net.NewSOCKS5DialContext(proxyURL)
		if err != nil {
			return nil, errors.New("grpc: could not create SOCKS5 dial context")
		}

		ctxDialer := func(ctx context.Context, addr string) (gonet.Conn, error) {
			return dialCtx(ctx, "tcp", addr)
		}

		grpcOpts = append(grpcOpts, grpc.WithContextDialer(ctxDialer))
	}

	keepAlive := keepalive.ClientParameters{
		PermitWithoutStream: opts.keepAlivePermitWithoutStream,
		Time:                MinKeepAliveTime,
		Timeout:             MinKeepAliveTimeout,
	}

	if opts.minKeepAliveTime > MinKeepAliveTime {
		keepAlive.Time = opts.minKeepAliveTime
	}

	if opts.minKeepAliveTimeout > MinKeepAliveTimeout {
		keepAlive.Timeout = opts.minKeepAliveTimeout
	}

	grpcOpts = append(grpcOpts, grpc.WithKeepaliveParams(keepAlive))

	if opts.tlsEnabled {
		fi, err := os.Stat(opts.tlsCertFilePath)
		if err != nil {
			return nil, err
		}

		var tlsCredentials credentials.TransportCredentials
		if fi.Mode().IsRegular() {
			tlsCredentials, err = credentials.NewClientTLSFromFile(opts.tlsCertFilePath, "")
			if err != nil {
				return nil, err
			}
		} else {
			tlsCredentials = credentials.NewTLS(&tls.Config{})
		}
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	return grpc.DialContext(context.Background(), target, grpcOpts...)
}

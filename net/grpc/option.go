package grpc

import (
	"time"

	"github.com/transientvariable/anchor"
)

// Option is a container for options used for configuring a gRPC connection.
type Option struct {
	keepAlivePermitWithoutStream bool
	messageSizeMaxReceive        int64
	messageSizeMaxSend           int64
	minKeepAliveTime             time.Duration
	minKeepAliveTimeout          time.Duration
	socks5Enabled                bool
	tlsCertFilePath              string
	tlsEnabled                   bool
	tlsKeyFilePath               string
}

// String returns a string representation of the Option.
func (o *Option) String() string {
	options := make(map[string]any)
	options["keep_alive_permit_without_stream"] = o.keepAlivePermitWithoutStream
	options["message_size_max_receive"] = o.messageSizeMaxReceive
	options["message_size_max_send"] = o.messageSizeMaxSend
	options["min_keep_alive_time"] = o.minKeepAliveTime
	options["min_keep_alive_timeout"] = o.minKeepAliveTimeout
	options["socks5_enabled"] = o.socks5Enabled
	options["tls_cert_file_path"] = o.tlsCertFilePath
	options["tls_enabled"] = o.tlsEnabled
	options["tls_key_file_path"] = o.tlsKeyFilePath
	return string(anchor.ToJSONFormatted(options))
}

// WithKeepAlivePermitWithoutStream ...
func WithKeepAlivePermitWithoutStream(keepAlive bool) func(*Option) {
	return func(o *Option) {
		o.keepAlivePermitWithoutStream = keepAlive
	}
}

// WithMessageSizeMaxReceive ...
func WithMessageSizeMaxReceive(size int64) func(*Option) {
	return func(o *Option) {
		o.messageSizeMaxReceive = size
	}
}

// WithMessageSizeMaxSend ...
func WithMessageSizeMaxSend(size int64) func(*Option) {
	return func(o *Option) {
		o.messageSizeMaxSend = size
	}
}

// WithMinKeepAliveTime ...
func WithMinKeepAliveTime(duration time.Duration) func(*Option) {
	return func(o *Option) {
		o.minKeepAliveTime = duration
	}
}

// WithMaxKeepAliveTimeout ...
func WithMaxKeepAliveTimeout(duration time.Duration) func(*Option) {
	return func(o *Option) {
		o.minKeepAliveTimeout = duration
	}
}

// WithSOCKS5Enabled ...
func WithSOCKS5Enabled(enable bool) func(*Option) {
	return func(o *Option) {
		o.socks5Enabled = enable
	}
}

// WithTLSCertFilePath ...
func WithTLSCertFilePath(path string) func(*Option) {
	return func(o *Option) {
		o.tlsCertFilePath = path
	}
}

// WithTLSKeyFilePath ...
func WithTLSKeyFilePath(path string) func(*Option) {
	return func(o *Option) {
		o.tlsKeyFilePath = path
	}
}

// WithTLSEnabled ...
func WithTLSEnabled(enable bool) func(*Option) {
	return func(o *Option) {
		o.tlsEnabled = enable
	}
}

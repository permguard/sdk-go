// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"
)

const defaultTimeout = 5 * time.Second

type clientOptions struct {
	timeout    time.Duration
	httpClient *http.Client
	tlsConfig  *tls.Config
	headers    http.Header
}

func defaultClientOptions() clientOptions {
	return clientOptions{timeout: defaultTimeout, headers: make(http.Header)}
}

// Option configures a Client.
type Option func(*clientOptions) error

// WithTimeout sets the per-call timeout used when the supplied context has no deadline.
func WithTimeout(timeout time.Duration) Option {
	return func(options *clientOptions) error {
		if timeout <= 0 {
			return errors.New("timeout must be greater than zero")
		}
		options.timeout = timeout
		return nil
	}
}

// WithHTTPClient supplies the HTTP client used by HTTP and HTTPS endpoints.
func WithHTTPClient(client *http.Client) Option {
	return func(options *clientOptions) error {
		if client == nil {
			return errors.New("HTTP client cannot be nil")
		}
		options.httpClient = client
		return nil
	}
}

// WithTLSConfig supplies TLS settings for HTTPS and gRPCS endpoints.
func WithTLSConfig(config *tls.Config) Option {
	return func(options *clientOptions) error {
		if config == nil {
			return errors.New("TLS config cannot be nil")
		}
		options.tlsConfig = config.Clone()
		return nil
	}
}

// WithHeader adds static HTTP headers or gRPC metadata to every request.
func WithHeader(name, value string) Option {
	return func(options *clientOptions) error {
		if name == "" {
			return errors.New("header name cannot be empty")
		}
		options.headers.Add(name, value)
		return nil
	}
}

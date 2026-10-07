// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	evaluationPath    = "/access/v1/evaluation"
	evaluationsPath   = "/access/v1/evaluations"
	configurationPath = "/.well-known/permguard-pdp-v1-configuration"
)

type transport interface {
	evaluate(context.Context, *EvaluateRequest, bool) (*EvaluateResponse, error)
	configuration(context.Context) (*Configuration, error)
	close() error
}

// Client evaluates access requests against a Permguard PDP over HTTP or gRPC.
type Client struct {
	transport transport
	timeout   time.Duration
}

// NewClient connects to an http(s):// or grpc(s):// PDP endpoint.
func NewClient(endpoint string, supplied ...Option) (*Client, error) {
	options := defaultClientOptions()
	for _, apply := range supplied {
		if apply == nil {
			continue
		}
		if err := apply(&options); err != nil {
			return nil, fmt.Errorf("configure Permguard client: %w", err)
		}
	}

	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse Permguard endpoint: %w", err)
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("parse Permguard endpoint: host is required")
	}

	var selected transport
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https":
		selected, err = newHTTPTransport(parsed, options)
	case "grpc", "grpcs":
		selected, err = newGRPCTransport(parsed, options)
	default:
		return nil, fmt.Errorf("unsupported Permguard endpoint scheme %q", parsed.Scheme)
	}
	if err != nil {
		return nil, err
	}

	return &Client{transport: selected, timeout: options.timeout}, nil
}

// Evaluate asks the single-evaluation endpoint for a decision.
func (c *Client) Evaluate(ctx context.Context, request *EvaluateRequest) (*EvaluateResponse, error) {
	return c.evaluate(ctx, request, false)
}

// EvaluateMany asks the boxcarred endpoint for decisions.
func (c *Client) EvaluateMany(ctx context.Context, request *EvaluateRequest) (*EvaluateResponse, error) {
	return c.evaluate(ctx, request, true)
}

func (c *Client) evaluate(ctx context.Context, request *EvaluateRequest, many bool) (*EvaluateResponse, error) {
	if c == nil || c.transport == nil {
		return nil, fmt.Errorf("Permguard client is not initialized")
	}
	if request == nil {
		return nil, fmt.Errorf("Permguard evaluation request cannot be nil")
	}
	ctx, cancel := withTimeout(ctx, c.timeout)
	defer cancel()
	return c.transport.evaluate(ctx, request, many)
}

// GetConfiguration returns the native PDP v1 discovery document.
func (c *Client) GetConfiguration(ctx context.Context) (*Configuration, error) {
	if c == nil || c.transport == nil {
		return nil, fmt.Errorf("Permguard client is not initialized")
	}
	ctx, cancel := withTimeout(ctx, c.timeout)
	defer cancel()
	return c.transport.configuration(ctx)
}

// Close releases transport resources. It is safe to call more than once.
func (c *Client) Close() error {
	if c == nil || c.transport == nil {
		return nil
	}
	return c.transport.close()
}

func withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, set := ctx.Deadline(); set {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, timeout)
}

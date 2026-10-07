// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const maxResponseBytes = 16 << 20

type httpTransport struct {
	base    string
	client  *http.Client
	headers http.Header
}

func newHTTPTransport(endpoint *url.URL, options clientOptions) (*httpTransport, error) {
	if endpoint.Path != "" && endpoint.Path != "/" {
		return nil, fmt.Errorf("HTTP Permguard endpoint must not contain a path")
	}
	if endpoint.RawQuery != "" || endpoint.Fragment != "" {
		return nil, fmt.Errorf("HTTP Permguard endpoint must not contain a query or fragment")
	}

	client := options.httpClient
	if client == nil {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		if options.tlsConfig != nil {
			transport.TLSClientConfig = options.tlsConfig.Clone()
		}
		client = &http.Client{Transport: transport}
	} else if options.tlsConfig != nil {
		return nil, fmt.Errorf("WithTLSConfig cannot be combined with WithHTTPClient; configure TLS on the supplied HTTP client")
	}

	return &httpTransport{
		base:    strings.TrimRight(endpoint.Scheme+"://"+endpoint.Host, "/"),
		client:  client,
		headers: options.headers.Clone(),
	}, nil
}

func (transport *httpTransport) evaluate(
	ctx context.Context,
	request *EvaluateRequest,
	many bool,
) (*EvaluateResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("encode Permguard evaluation request: %w", err)
	}
	path := evaluationPath
	if many {
		path = evaluationsPath
	}

	var response EvaluateResponse
	if err := transport.call(ctx, http.MethodPost, path, body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (transport *httpTransport) configuration(ctx context.Context) (*Configuration, error) {
	var configuration Configuration
	if err := transport.call(ctx, http.MethodGet, configurationPath, nil, &configuration); err != nil {
		return nil, err
	}
	return &configuration, nil
}

func (transport *httpTransport) call(
	ctx context.Context,
	method string,
	path string,
	body []byte,
	target any,
) error {
	request, err := http.NewRequestWithContext(ctx, method, transport.base+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create Permguard HTTP request: %w", err)
	}
	request.Header = transport.headers.Clone()
	request.Header.Set("Accept", "application/json")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := transport.client.Do(request)
	if err != nil {
		return fmt.Errorf("call Permguard HTTP endpoint: %w", err)
	}
	defer response.Body.Close()

	payload, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil {
		return fmt.Errorf("read Permguard HTTP response: %w", err)
	}
	if len(payload) > maxResponseBytes {
		return fmt.Errorf("read Permguard HTTP response: body exceeds %d bytes", maxResponseBytes)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return refusalFromHTTP(response.StatusCode, payload)
	}
	if err := json.Unmarshal(payload, target); err != nil {
		return fmt.Errorf("decode Permguard HTTP response: %w", err)
	}
	return nil
}

func refusalFromHTTP(status int, payload []byte) error {
	refusal := &Refusal{HTTPStatus: status}
	if err := json.Unmarshal(payload, refusal); err != nil || refusal.Message == "" {
		refusal.Message = strings.TrimSpace(string(payload))
		if refusal.Message == "" {
			refusal.Message = http.StatusText(status)
		}
	}
	if refusal.Class == "" {
		refusal.Class = httpClass(status)
	}
	if refusal.Code == "" {
		refusal.Code = "http_status"
	}
	return refusal
}

func httpClass(status int) string {
	switch status {
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		return "validation"
	case http.StatusConflict:
		return "conflict"
	case http.StatusUnauthorized, http.StatusForbidden:
		return "authorization"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return "unavailable"
	default:
		return "internal"
	}
}

func (transport *httpTransport) close() error {
	transport.client.CloseIdleConnections()
	return nil
}

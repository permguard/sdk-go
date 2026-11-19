// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	pdpv1 "github.com/permguard/sdk-go/internal/grpc/v1"
	"google.golang.org/grpc"
)

// silentPDP accepts every call and answers none until its test ends.
type silentPDP struct {
	pdpv1.UnimplementedPolicyDecisionPointServer
	done chan struct{}
}

func (server *silentPDP) Evaluate(ctx context.Context, _ *pdpv1.EvaluateRequest) (*pdpv1.EvaluateResponse, error) {
	select {
	case <-ctx.Done():
	case <-server.done:
	}
	return nil, ctx.Err()
}

// silentEndpoints serves a data plane that never answers, over HTTP and gRPC.
func silentEndpoints(t *testing.T) (string, string) {
	t.Helper()
	done := make(chan struct{})
	httpServer := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		select {
		case <-request.Context().Done():
		case <-done:
		}
	}))
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pdpv1.RegisterPolicyDecisionPointServer(grpcServer, &silentPDP{done: done})
	go func() { _ = grpcServer.Serve(listener) }()
	t.Cleanup(func() {
		close(done)
		grpcServer.Stop()
		httpServer.Close()
	})
	return httpServer.URL, "grpc://" + listener.Addr().String()
}

func refusalOf(t *testing.T, err error) *Refusal {
	t.Helper()
	var refusal *Refusal
	if !errors.As(err, &refusal) {
		t.Fatalf("expected a Refusal, got %T: %v", err, err)
	}
	return refusal
}

// BUG-0310: WithTimeout bounds every call, even when the context's deadline is later; a context
// deadline that comes sooner still ends the call first. BUG-0284: a timeout is the same refusal on
// both transports, class unavailable and code timeout.
func TestTheTimeoutBoundsEveryCallOnBothTransports(t *testing.T) {
	httpURL, grpcURL := silentEndpoints(t)
	for _, endpoint := range []string{httpURL, grpcURL} {
		client, err := NewClient(endpoint, WithTimeout(300*time.Millisecond))
		if err != nil {
			t.Fatal(err)
		}
		request := &EvaluateRequest{Zone: "acme", Ledger: "main"}

		long, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		started := time.Now()
		_, err = client.Evaluate(long, request)
		cancel()
		if elapsed := time.Since(started); elapsed > 2*time.Second {
			t.Fatalf("%s: the client's timeout did not bound the call: %v", endpoint, elapsed)
		}
		refusal := refusalOf(t, err)
		if refusal.Class != ClassUnavailable || refusal.Code != CodeTimeout {
			t.Fatalf("%s: timeout refusal = %+v", endpoint, refusal)
		}

		short, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		started = time.Now()
		_, err = client.Evaluate(short, request)
		cancel()
		if elapsed := time.Since(started); elapsed > 250*time.Millisecond {
			t.Fatalf("%s: the context's sooner deadline did not end the call: %v", endpoint, elapsed)
		}
		if refusal := refusalOf(t, err); refusal.Code != CodeTimeout {
			t.Fatalf("%s: %+v", endpoint, refusal)
		}

		// The caller's own cancellation is not a refusal.
		cancelled, cancel := context.WithCancel(context.Background())
		go func() { time.Sleep(50 * time.Millisecond); cancel() }()
		_, err = client.Evaluate(cancelled, request)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("%s: a cancellation is the context's error, not %T: %v", endpoint, err, err)
		}
		_ = client.Close()
	}
}

// BUG-0284: a data plane that cannot be reached is the same refusal on both transports.
func TestAnUnreachableDataPlaneIsTheSameRefusalOnBothTransports(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	_ = listener.Close()
	for _, endpoint := range []string{"http://" + address, "grpc://" + address} {
		client, err := NewClient(endpoint, WithTimeout(2*time.Second))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Evaluate(context.Background(), &EvaluateRequest{Zone: "acme", Ledger: "main"})
		refusal := refusalOf(t, err)
		if refusal.Class != ClassUnavailable || refusal.Code != CodeUnreachable {
			t.Fatalf("%s: unreachable refusal = %+v", endpoint, refusal)
		}
		_ = client.Close()
	}
}

// BUG-0317: a gRPC endpoint without a port dials the port its scheme implies, 80 in plaintext and
// 443 over TLS; and a client closes more than once.
func TestAGRPCEndpointWithoutAPortAndASecondClose(t *testing.T) {
	for endpoint, target := range map[string]string{
		"grpc://127.0.0.1":  "127.0.0.1:80",
		"grpcs://127.0.0.1": "127.0.0.1:443",
		"grpc://[::1]":      "[::1]:80",
	} {
		client, err := NewClient(endpoint)
		if err != nil {
			t.Fatal(err)
		}
		if got := client.transport.(*grpcTransport).connection.CanonicalTarget(); got != "dns:///"+target && got != "passthrough:///"+target {
			t.Fatalf("%s dials %s, not %s", endpoint, got, target)
		}
		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
		if err := client.Close(); err != nil {
			t.Fatalf("%s: second Close: %v", endpoint, err)
		}
	}
}

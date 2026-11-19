// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"

	pdpv1 "github.com/permguard/sdk-go/internal/grpc/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	grpcErrorClass = "permguard-error-class"
	grpcErrorCode  = "permguard-error-code"
)

type grpcTransport struct {
	closing    sync.Once
	closed     error
	connection *grpc.ClientConn
	client     pdpv1.PolicyDecisionPointClient
	metadata   metadata.MD
}

func newGRPCTransport(endpoint *url.URL, options clientOptions) (*grpcTransport, error) {
	if endpoint.Path != "" && endpoint.Path != "/" {
		return nil, fmt.Errorf("gRPC Permguard endpoint must not contain a path")
	}
	if endpoint.RawQuery != "" || endpoint.Fragment != "" {
		return nil, fmt.Errorf("gRPC Permguard endpoint must not contain a query or fragment")
	}

	var transportCredentials credentials.TransportCredentials
	if strings.EqualFold(endpoint.Scheme, "grpcs") {
		configuration := options.tlsConfig
		if configuration == nil {
			configuration = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		transportCredentials = credentials.NewTLS(configuration.Clone())
	} else {
		transportCredentials = insecure.NewCredentials()
	}

	// The port the scheme implies when the endpoint names none, as the Permguard CLI's endpoint
	// grammar has it: 443 over TLS, 80 in plaintext, never the TLS port for plaintext (BUG-0317).
	target := endpoint.Host
	if endpoint.Port() == "" {
		port := "80"
		if strings.EqualFold(endpoint.Scheme, "grpcs") {
			port = "443"
		}
		target = net.JoinHostPort(endpoint.Hostname(), port)
	}
	connection, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(transportCredentials),
	)
	if err != nil {
		return nil, fmt.Errorf("create Permguard gRPC client: %w", err)
	}

	return &grpcTransport{
		connection: connection,
		client:     pdpv1.NewPolicyDecisionPointClient(connection),
		metadata:   metadataFromHeaders(options.headers),
	}, nil
}

func metadataFromHeaders(headers http.Header) metadata.MD {
	values := metadata.MD{}
	for name, held := range headers {
		key := strings.ToLower(name)
		values[key] = append([]string(nil), held...)
	}
	return values
}

func (transport *grpcTransport) evaluate(
	ctx context.Context,
	request *EvaluateRequest,
	many bool,
) (*EvaluateResponse, error) {
	wire, err := requestToProto(request)
	if err != nil {
		return nil, fmt.Errorf("encode Permguard gRPC request: %w", err)
	}
	ctx = transport.outgoing(ctx)

	var header metadata.MD
	var trailer metadata.MD
	var response *pdpv1.EvaluateResponse
	if many {
		response, err = transport.client.EvaluateMany(
			ctx,
			wire,
			grpc.Header(&header),
			grpc.Trailer(&trailer),
		)
	} else {
		response, err = transport.client.Evaluate(
			ctx,
			wire,
			grpc.Header(&header),
			grpc.Trailer(&trailer),
		)
	}
	if err != nil {
		return nil, refusalFromGRPC(ctx, err, header, trailer)
	}
	return responseFromProto(response), nil
}

func (transport *grpcTransport) configuration(ctx context.Context) (*Configuration, error) {
	ctx = transport.outgoing(ctx)
	var header metadata.MD
	var trailer metadata.MD
	response, err := transport.client.GetConfiguration(
		ctx,
		&pdpv1.GetConfigurationRequest{},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		return nil, refusalFromGRPC(ctx, err, header, trailer)
	}
	return configurationFromProto(response), nil
}

func (transport *grpcTransport) outgoing(ctx context.Context) context.Context {
	if len(transport.metadata) == 0 {
		return ctx
	}
	existing, _ := metadata.FromOutgoingContext(ctx)
	return metadata.NewOutgoingContext(ctx, metadata.Join(existing, transport.metadata))
}

func refusalFromGRPC(ctx context.Context, err error, header, trailer metadata.MD) error {
	held, ok := status.FromError(err)
	if !ok {
		if failure := clientFailure(ctx, err); failure != nil {
			return failure
		}
		return fmt.Errorf("call Permguard gRPC endpoint: %w", err)
	}
	lookup := func(key string) string {
		if value := trailer.Get(key); len(value) > 0 {
			return value[0]
		}
		if value := header.Get(key); len(value) > 0 {
			return value[0]
		}
		return ""
	}
	// What the client's own gRPC library reported, not the data plane: a timeout, the caller's
	// cancellation, or a data plane it could not reach, as the HTTP transport reports them
	// (BUG-0284). A status the data plane sent carries its own class and code.
	if lookup(grpcErrorCode) == "" {
		switch held.Code() {
		case codes.Canceled:
			if ctx.Err() != nil {
				return ctx.Err()
			}
		case codes.DeadlineExceeded:
			return &Refusal{
				Class:    ClassUnavailable,
				Code:     CodeTimeout,
				Message:  "the data plane did not answer within the call's timeout",
				GRPCCode: held.Code(),
			}
		case codes.Unavailable:
			return &Refusal{
				Class:    ClassUnavailable,
				Code:     CodeUnreachable,
				Message:  "the data plane could not be reached: " + held.Message(),
				GRPCCode: held.Code(),
			}
		}
	}
	class := lookup(grpcErrorClass)
	if class == "" {
		class = grpcClass(held.Code())
	}
	code := lookup(grpcErrorCode)
	if code == "" {
		code = strings.ToLower(held.Code().String())
	}
	return &Refusal{
		Class:    class,
		Code:     code,
		Message:  held.Message(),
		GRPCCode: held.Code(),
	}
}

func grpcClass(code codes.Code) string {
	switch code {
	case codes.InvalidArgument, codes.OutOfRange:
		return "validation"
	case codes.FailedPrecondition, codes.AlreadyExists, codes.Aborted:
		return "conflict"
	case codes.Unauthenticated, codes.PermissionDenied:
		return "authorization"
	case codes.NotFound:
		return "not_found"
	case codes.Unavailable, codes.DeadlineExceeded:
		return "unavailable"
	default:
		return "internal"
	}
}

// close is safe to call more than once: a connection already closed is not an error (BUG-0317).
func (transport *grpcTransport) close() error {
	transport.closing.Do(func() {
		transport.closed = transport.connection.Close()
	})
	return transport.closed
}

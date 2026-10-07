// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

	connection, err := grpc.NewClient(
		endpoint.Host,
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
		return nil, refusalFromGRPC(err, header, trailer)
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
		return nil, refusalFromGRPC(err, header, trailer)
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

func refusalFromGRPC(err error, header, trailer metadata.MD) error {
	held, ok := status.FromError(err)
	if !ok {
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

func (transport *grpcTransport) close() error {
	return transport.connection.Close()
}

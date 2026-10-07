// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"context"
	"errors"
	"net"
	"testing"

	pdpv1 "github.com/permguard/sdk-go/internal/grpc/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type testPDPServer struct {
	testing *testing.T
}

func (server *testPDPServer) Evaluate(
	ctx context.Context,
	request *pdpv1.EvaluateRequest,
) (*pdpv1.EvaluateResponse, error) {
	server.assertMetadata(ctx)
	if request.Ledger == "bad" {
		_ = grpc.SetTrailer(ctx, metadata.Pairs(
			grpcErrorClass, "validation",
			grpcErrorCode, "ledger_invalid",
		))
		return nil, status.Error(codes.InvalidArgument, "ledger is invalid")
	}
	return &pdpv1.EvaluateResponse{
		Decision:  true,
		RequestId: request.RequestId,
		Context: &pdpv1.DecisionContext{
			Policies: []string{"policy-1"},
		},
	}, nil
}

func (server *testPDPServer) EvaluateMany(
	ctx context.Context,
	request *pdpv1.EvaluateRequest,
) (*pdpv1.EvaluateResponse, error) {
	server.assertMetadata(ctx)
	return &pdpv1.EvaluateResponse{
		Decision: false,
		Evaluations: []*pdpv1.Decision{
			{Decision: true, RequestId: "one"},
			{Decision: false, RequestId: "two"},
		},
	}, nil
}

func (server *testPDPServer) GetConfiguration(
	ctx context.Context,
	_ *pdpv1.GetConfigurationRequest,
) (*pdpv1.GetConfigurationResponse, error) {
	server.assertMetadata(ctx)
	return &pdpv1.GetConfigurationResponse{
		Interface: "permguard.api.pdp.native.v1",
		Pdp:       "grpc://test",
	}, nil
}

func (server *testPDPServer) assertMetadata(ctx context.Context) {
	held, _ := metadata.FromIncomingContext(ctx)
	if got := held.Get("authorization"); len(got) != 1 || got[0] != "Bearer test" {
		server.testing.Errorf("authorization metadata = %#v", got)
	}
}

func TestGRPCTransportServesNativePDPContract(t *testing.T) {
	if pdpv1.PolicyDecisionPoint_Evaluate_FullMethodName != "/permguard.data.v1.PolicyDecisionPoint/Evaluate" ||
		pdpv1.PolicyDecisionPoint_EvaluateMany_FullMethodName != "/permguard.data.v1.PolicyDecisionPoint/EvaluateMany" ||
		pdpv1.PolicyDecisionPoint_GetConfiguration_FullMethodName != "/permguard.data.v1.PolicyDecisionPoint/GetConfiguration" {
		t.Fatal("generated gRPC method names do not match the native PDP v1 contract")
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := grpc.NewServer()
	pdpv1.RegisterPolicyDecisionPointServer(server, &testPDPServer{testing: t})
	go func() {
		_ = server.Serve(listener)
	}()
	defer server.Stop()

	client, err := NewClient("grpc://"+listener.Addr().String(), WithHeader("Authorization", "Bearer test"))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	defer client.Close()

	decision, err := client.Evaluate(context.Background(), &EvaluateRequest{
		Zone: "acme", Ledger: "documents", RequestID: "r1",
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !decision.Decision || decision.RequestID != "r1" {
		t.Fatalf("decision = %#v", decision)
	}

	batch, err := client.EvaluateMany(context.Background(), &EvaluateRequest{
		Zone: "acme", Ledger: "documents", Evaluations: []Evaluation{{}, {}},
	})
	if err != nil {
		t.Fatalf("evaluate many: %v", err)
	}
	if batch.Decision || len(batch.Evaluations) != 2 {
		t.Fatalf("batch = %#v", batch)
	}

	configuration, err := client.GetConfiguration(context.Background())
	if err != nil {
		t.Fatalf("configuration: %v", err)
	}
	if configuration.Interface != "permguard.api.pdp.native.v1" {
		t.Fatalf("interface = %q", configuration.Interface)
	}

	_, err = client.Evaluate(context.Background(), &EvaluateRequest{Zone: "acme", Ledger: "bad"})
	var refusal *Refusal
	if !errors.As(err, &refusal) {
		t.Fatalf("error = %T %v, want *Refusal", err, err)
	}
	if refusal.Code != "ledger_invalid" || refusal.Class != "validation" || refusal.GRPCCode != codes.InvalidArgument {
		t.Fatalf("refusal = %#v", refusal)
	}
}

func TestGRPCMapperRejectsLossyIntegers(t *testing.T) {
	_, err := requestToProto(&EvaluateRequest{
		Zone:    "acme",
		Ledger:  "documents",
		Context: map[string]any{"too_large": int64(maxExactProtoInteger + 1)},
	})
	if err == nil {
		t.Fatal("integer larger than 2^53 was silently accepted")
	}
}

func TestGRPCConflictFallbackMatchesSharedContract(t *testing.T) {
	for _, code := range []codes.Code{codes.FailedPrecondition, codes.AlreadyExists, codes.Aborted} {
		if got := grpcClass(code); got != "conflict" {
			t.Fatalf("gRPC %s class = %q, want conflict", code, got)
		}
	}
}

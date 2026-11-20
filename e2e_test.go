// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

//go:build e2e

package permguard

import (
	"context"
	"errors"
	"os"
	"testing"
)

// The data plane `scripts/e2e-server.sh` started, from the Permguard release PERMGUARD_CONTRACT
// pins, serving a ledger `e2e/policies` that permits `User:alice` to `read` (BUG-0077).
func e2eEndpoints(t *testing.T) map[string]string {
	t.Helper()
	http, grpc := os.Getenv("E2E_HTTP"), os.Getenv("E2E_GRPC")
	if http == "" || grpc == "" {
		t.Fatal("E2E_HTTP and E2E_GRPC name the data plane scripts/e2e-server.sh started")
	}
	return map[string]string{"HTTP": http, "gRPC": grpc}
}

func ask(id string) *EvaluateRequest {
	return &EvaluateRequest{
		Zone:     "e2e",
		Ledger:   "policies",
		Subject:  &Entity{Type: "User", ID: id},
		Resource: &Entity{Type: "Document", ID: "d"},
		Action:   &Action{Name: "read"},
	}
}

// Against a real data plane, over both transports: a permit with its decision evidence, a deny, a
// batch, a refusal, a request for a signed response, and the discovery document (BUG-0077,
// BUG-0274).
func TestAgainstARealDataPlane(t *testing.T) {
	for name, endpoint := range e2eEndpoints(t) {
		t.Run(name, func(t *testing.T) {
			client, err := NewClient(endpoint)
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			ctx := context.Background()

			permit, err := client.Evaluate(ctx, ask("alice"))
			if err != nil || !permit.Decision {
				t.Fatalf("alice reads: %+v %v", permit, err)
			}
			if permit.Context == nil || permit.Context.Evidence == nil ||
				permit.Context.Evidence.Status == "" {
				t.Fatalf("the decision's evidence: %+v", permit.Context)
			}
			deny, err := client.Evaluate(ctx, ask("bob"))
			if err != nil || deny.Decision {
				t.Fatalf("bob reads: %+v %v", deny, err)
			}

			batch := ask("alice")
			batch.Evaluations = []Evaluation{
				{Subject: &Entity{Type: "User", ID: "alice"}},
				{Subject: &Entity{Type: "User", ID: "bob"}},
			}
			many, err := client.EvaluateMany(ctx, batch)
			if err != nil || len(many.Evaluations) != 2 ||
				!many.Evaluations[0].Decision || many.Evaluations[1].Decision {
				t.Fatalf("the batch: %+v %v", many, err)
			}

			missing := ask("alice")
			missing.Ledger = "nope"
			_, err = client.Evaluate(ctx, missing)
			var refusal *Refusal
			if !errors.As(err, &refusal) || refusal.Code != "ledger_not_served" {
				t.Fatalf("an unknown ledger: %v", err)
			}

			// The nonce reaches the plane: a signed response is addressed to an authenticated
			// caller, and this one presents no token, so the plane refuses rather than answer
			// unsigned. The attestation's mapping is the unit tests'.
			signed := ask("alice")
			signed.RequestNonce = "ZTJlLW5vbmNlLTE2LWJ5dGVz"
			_, err = client.Evaluate(ctx, signed)
			if !errors.As(err, &refusal) || refusal.Code != "unauthenticated" {
				t.Fatalf("a nonce without a credential: %v", err)
			}

			configuration, err := client.GetConfiguration(ctx)
			if err != nil || configuration.Version == "" || configuration.Stage == "" ||
				configuration.Limits["max_evaluations"] == 0 || configuration.ResponseSigning.State == "" {
				t.Fatalf("the discovery document: %+v %v", configuration, err)
			}
		})
	}
}

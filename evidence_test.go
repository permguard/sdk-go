// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"encoding/json"
	"strings"
	"testing"

	pdpv1 "github.com/permguard/sdk-go/internal/grpc/v1"
)

// A request asks for a signed response with its nonce, and the answer's attestation and decision
// evidence reach the caller over both transports (BUG-0274).
func TestTheNonceTheAttestationAndTheEvidenceTravelOnBothTransports(t *testing.T) {
	request := &EvaluateRequest{Zone: "z", Ledger: "l", RequestNonce: "bm9uY2Utb2Ytc2l4dGVlbg"}
	mapped, err := requestToProto(request)
	if err != nil {
		t.Fatal(err)
	}
	if mapped.RequestNonce != request.RequestNonce {
		t.Fatalf("gRPC nonce: %q", mapped.RequestNonce)
	}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(encoded) || !strings.Contains(string(encoded), `"request_nonce":"bm9uY2Utb2Ytc2l4dGVlbg"`) {
		t.Fatalf("HTTP nonce: %s", encoded)
	}

	record := "rec-1"
	fromGRPC := responseFromProto(&pdpv1.EvaluateResponse{
		Decision: true,
		Context: &pdpv1.DecisionContext{
			Id: "d-1",
			Evidence: &pdpv1.Evidence{
				Mode:     "required",
				Status:   "recorded",
				RecordId: &record,
				Position: &pdpv1.EvidencePosition{HostId: "h", Run: "r", Sequence: "7"},
			},
		},
		Attestation: &pdpv1.Attestation{Protected: "p", Payload: "b", Signature: "s"},
	})
	var fromHTTP EvaluateResponse
	if err := json.Unmarshal([]byte(`{"decision":true,"context":{"id":"d-1","evidence":{"mode":"required","status":"recorded","record_id":"rec-1","position":{"host_id":"h","run":"r","sequence":"7"}}},"attestation":{"protected":"p","payload":"b","signature":"s"}}`), &fromHTTP); err != nil {
		t.Fatal(err)
	}
	for name, held := range map[string]*EvaluateResponse{"gRPC": fromGRPC, "HTTP": &fromHTTP} {
		if held.Attestation == nil || held.Attestation.Signature != "s" {
			t.Fatalf("%s: attestation %+v", name, held.Attestation)
		}
		evidence := held.Context.Evidence
		if evidence == nil || evidence.Status != "recorded" || evidence.RecordID == nil ||
			*evidence.RecordID != "rec-1" || evidence.Position == nil || evidence.Position.Sequence != "7" {
			t.Fatalf("%s: evidence %+v", name, evidence)
		}
	}
}

// The discovery document's version, lifecycle, limits and response signing reach the caller over
// both transports (BUG-0274).
func TestTheDiscoveryDocumentIsReadWholeOnBothTransports(t *testing.T) {
	fromGRPC := configurationFromProto(&pdpv1.GetConfigurationResponse{
		Interface: "permguard.api.pdp.native.v1",
		Version:   "0.1.33",
		Component: "data-plane",
		State:     "serving",
		Stage:     "v1",
		Limits:    map[string]uint64{"max_evaluations": 64},
		ResponseSigning: &pdpv1.ResponseSigning{
			State: "optional", ValiditySeconds: 300, ClockSkewSeconds: 30,
		},
	})
	var fromHTTP Configuration
	if err := json.Unmarshal([]byte(`{"interface":"permguard.api.pdp.native.v1","version":"0.1.33","component":"data-plane","state":"serving","stage":"v1","limits":{"max_evaluations":64},"response_signing":{"state":"optional","validity_seconds":300,"clock_skew_seconds":30}}`), &fromHTTP); err != nil {
		t.Fatal(err)
	}
	for name, held := range map[string]*Configuration{"gRPC": fromGRPC, "HTTP": &fromHTTP} {
		if held.Version != "0.1.33" || held.Component != "data-plane" || held.State != "serving" ||
			held.Stage != "v1" || held.Limits["max_evaluations"] != 64 ||
			held.ResponseSigning.State != "optional" || held.ResponseSigning.ValiditySeconds != 300 {
			t.Fatalf("%s: %+v", name, held)
		}
	}
}

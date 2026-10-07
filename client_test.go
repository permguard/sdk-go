// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPTransportServesNativePDPContract(t *testing.T) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case evaluationPath:
			if request.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", request.Method)
			}
			if got := request.Header.Get("Authorization"); got != "Bearer test" {
				t.Errorf("authorization = %q", got)
			}
			var payload EvaluateRequest
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Errorf("decode request: %v", err)
			}
			if payload.Zone != "acme" || payload.Ledger != "documents" {
				t.Errorf("store = %s/%s", payload.Zone, payload.Ledger)
			}
			response.Header().Set("Content-Type", "application/json")
			_, _ = response.Write([]byte(`{"decision":false,"request_id":"r1","context":{"policies":["policy-1"]}}`))
		case configurationPath:
			response.Header().Set("Content-Type", "application/json")
			_, _ = response.Write([]byte(`{"interface":"permguard.api.pdp.native.v1","pdp":"test","endpoints":{"evaluation":"e","evaluations":"es"},"capabilities":[],"store_scope":{"in":"payload","zone":"required","ledger":"required","profile":"optional"}}`))
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL, WithHeader("Authorization", "Bearer test"))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	defer client.Close()

	decision, err := client.Evaluate(context.Background(), &EvaluateRequest{
		Zone:      "acme",
		Ledger:    "documents",
		RequestID: "r1",
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if decision.Decision {
		t.Fatal("deny was changed into permit")
	}
	if len(decision.Context.Policies) != 1 || decision.Context.Policies[0] != "policy-1" {
		t.Fatalf("policies = %#v", decision.Context.Policies)
	}

	configuration, err := client.GetConfiguration(context.Background())
	if err != nil {
		t.Fatalf("configuration: %v", err)
	}
	if configuration.Interface != "permguard.api.pdp.native.v1" {
		t.Fatalf("interface = %q", configuration.Interface)
	}
}

func TestHTTPRefusalRemainsStructured(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.WriteHeader(http.StatusBadRequest)
		_, _ = response.Write([]byte(`{"class":"validation","code":"zone_required","message":"zone is required"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	defer client.Close()

	_, err = client.Evaluate(context.Background(), &EvaluateRequest{})
	var refusal *Refusal
	if !errors.As(err, &refusal) {
		t.Fatalf("error = %T %v, want *Refusal", err, err)
	}
	if refusal.Code != "zone_required" || refusal.Class != "validation" || refusal.HTTPStatus != 400 {
		t.Fatalf("refusal = %#v", refusal)
	}
}

func TestEvaluationPartitionInputPresenceIsPreservedInJSON(t *testing.T) {
	empty := PartitionInputs{}
	emptyContext := map[string]any{}
	payload, err := json.Marshal(EvaluateRequest{
		Zone:        "acme",
		Ledger:      "documents",
		Evaluations: []Evaluation{{Context: &emptyContext, PartitionInputs: &empty}},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	evaluation := decoded["evaluations"].([]any)[0].(map[string]any)
	if _, present := evaluation["partition_inputs"]; !present {
		t.Fatalf("partition_inputs presence was lost: %s", payload)
	}
	if _, present := evaluation["context"]; !present {
		t.Fatalf("context presence was lost: %s", payload)
	}
}

func TestHTTPFallbackAndEndpointValidationMatchSharedContract(t *testing.T) {
	if got := httpClass(http.StatusConflict); got != "conflict" {
		t.Fatalf("HTTP 409 class = %q, want conflict", got)
	}
	if _, err := NewClient("http://user:secret@pdp.example"); err == nil {
		t.Fatal("endpoint with embedded credentials was accepted")
	}
}

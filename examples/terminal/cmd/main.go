// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	permguard "github.com/permguard/sdk-go"
)

func main() {
	endpoint := os.Getenv("PERMGUARD_PDP_URL")
	if endpoint == "" {
		endpoint = "grpc://localhost:7443"
	}

	client, err := permguard.NewClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	response, err := client.Evaluate(context.Background(), &permguard.EvaluateRequest{
		Zone:    "acme",
		Ledger:  "main-ledger",
		Profile: "gateway",
		Subject: &permguard.Entity{Type: "User", ID: "alice"},
		Resource: &permguard.Entity{
			Type: "Document",
			ID:   "budget-2026",
		},
		Action:    &permguard.Action{Name: "read"},
		RequestID: "example-1",
	})
	if err != nil {
		var refusal *permguard.Refusal
		if errors.As(err, &refusal) {
			log.Fatalf("PDP refused the request: class=%s code=%s message=%s", refusal.Class, refusal.Code, refusal.Message)
		}
		log.Fatal(err)
	}

	fmt.Printf("permitted: %t\n", response.Decision)
}

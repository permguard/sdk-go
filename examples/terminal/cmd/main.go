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
		Ledger:  "documents",
		Subject: &permguard.Entity{Type: "user", ID: "amy@example.com"},
		Resource: &permguard.Entity{
			Type: "document",
			ID:   "quarterly-report",
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

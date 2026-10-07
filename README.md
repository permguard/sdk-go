<!--
Copyright (c) 2022 Nitro Agility S.r.l.
SPDX-License-Identifier: Apache-2.0
-->

# Permguard Go SDK

The official Go client for the stateless Permguard PDP interface
`permguard.api.pdp.native.v1`.

One public API supports both server bindings:

- `http://` and `https://` use the JSON endpoints;
- `grpc://` and `grpcs://` use `permguard.data.v1.PolicyDecisionPoint`.

The request and response semantics are identical on both transports.

## Requirements

- Go 1.23.5 or newer.

## Installation

```bash
go get github.com/permguard/sdk-go
```

## Evaluate one request

```go
package main

import (
	"context"
	"fmt"
	"log"

	permguard "github.com/permguard/sdk-go"
)

func main() {
	client, err := permguard.NewClient("grpc://localhost:7443")
	// Use http://localhost:7443 for the HTTP/JSON binding.
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	response, err := client.Evaluate(context.Background(), &permguard.EvaluateRequest{
		Zone:   "acme",
		Ledger: "documents",
		Subject: &permguard.Entity{
			Type: "user",
			ID:   "amy@example.com",
		},
		Resource: &permguard.Entity{
			Type: "document",
			ID:   "quarterly-report",
		},
		Action: &permguard.Action{Name: "read"},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("permitted:", response.Decision)
}
```

A deny is a successful response with `Decision == false`. Validation,
authorization, availability, and server failures are returned as
`*permguard.Refusal`, preserving their stable class and code.

## Partition inputs

Runtime data is addressed to the partition name declared by the ledger profile:

```go
request.PartitionInputs = permguard.PartitionInputs{
	"authorization": {
		Type: "permguard.cedar.entities.v1",
		Data: []any{
			map[string]any{
				"uid": map[string]any{"type": "Team", "id": "engineering"},
				"attrs": map[string]any{"active": true},
				"parents": []any{},
			},
		},
	},
}
```

The map key must be the partition name from the selected profile. The `Type`
asserts the partition input contract; it does not select a policy runtime.

## Evaluate a batch

Set top-level defaults, add `Evaluations`, and call `EvaluateMany`:

```go
response, err := client.EvaluateMany(ctx, &permguard.EvaluateRequest{
	Zone:    "acme",
	Ledger:  "documents",
	Subject: &permguard.Entity{Type: "user", ID: "amy@example.com"},
	Evaluations: []permguard.Evaluation{
		{
			Resource:  &permguard.Entity{Type: "document", ID: "one"},
			Action:    &permguard.Action{Name: "read"},
			RequestID: "one",
		},
		{
			Resource:  &permguard.Entity{Type: "document", ID: "two"},
			Action:    &permguard.Action{Name: "read"},
			RequestID: "two",
		},
	},
	Options: &permguard.EvaluationOptions{
		EvaluationsSemantic: permguard.ExecuteAll,
	},
})
```

## Configuration and transport options

```go
configuration, err := client.GetConfiguration(ctx)
```

Available options include:

- `WithTimeout`
- `WithTLSConfig`
- `WithHTTPClient`
- `WithHeader`, also carried as gRPC metadata

The client reuses its HTTP transport or gRPC channel. Call `Close` when it is no
longer needed.

## Compatibility

This major version implements `permguard.api.pdp.native.v1`. Compatibility is
tied to that versioned interface rather than to an unrelated server minor
version.

## Development

```bash
make protoc
make test
make lint
```

`proto/v1/pdp.proto` is generated from the canonical contract in the Permguard
server repository. Changes to the contract must regenerate the Go bindings and
pass both HTTP and gRPC contract tests.

## License

Apache License 2.0. See [LICENSE](LICENSE) and
[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

<!--
Copyright (c) 2022 Nitro Agility S.r.l.
SPDX-License-Identifier: Apache-2.0
-->

# Permguard Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/permguard/sdk-go.svg)](https://pkg.go.dev/github.com/permguard/sdk-go)
[![licence](https://img.shields.io/badge/licence-Apache--2.0-blue)](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center">
  <img src="https://raw.githubusercontent.com/permguard/sdk-go/main/assets/readme/hero.png?v=3" alt="Permguard SDK for Go. Ask Permguard, from Go: the official Go client of the Permguard Policy Decision Point, Go 1.23.5 or newer, over HTTP or gRPC, Apache-2.0. Beside it, a Go window creates a client for grpc://localhost:7443 and evaluates whether the User alice may read in the ledger main-ledger of the zone acme; the answer is permit." width="820">
</p>

<p align="center">
  <a href="https://github.com/permguard/permguard"><b>Permguard</b></a> ·
  <a href="https://docs.permguard.com/"><b>Documentation</b></a> ·
  <a href="#install"><b>Install</b></a> ·
  <a href="#quick-start"><b>Quick start</b></a> ·
  <a href="#develop"><b>Develop</b></a>
</p>

**The official Go client of the Permguard Policy Decision Point.**

Before your code acts, it asks Permguard whether a subject may perform an action on a resource, and Permguard answers from the policies of a ledger.
This SDK is how Go code asks: one client for the native, stateless PDP interface `permguard.api.pdp.native.v1`, over HTTP/JSON or gRPC, with the same request and the same answer on both.

<p align="center">
  <img src="https://raw.githubusercontent.com/permguard/sdk-go/main/assets/readme/flow.png?v=2" alt="How a decision is asked. Your Go code asks before it acts: who, User alice; does what, read; on what, Document budget-2026. The Permguard Go SDK, one client with Evaluate and EvaluateMany, sends it as JSON to an http or https endpoint or over gRPC to a grpc or grpcs endpoint. The Permguard data plane, the PDP of the native v1 interface, decides with the policies of the ledger main-ledger in the zone acme, under the profile gateway, on Cedar, Rego or Dogwood. The answer is a permit, a deny, which is a normal answer, or a *permguard.Refusal with its class and code." width="820">
</p>

<p align="center"><sub>Your code asks, the SDK carries the question, the data plane decides with the ledger's policies.</sub></p>

## Install

The SDK needs Go 1.23.5 or newer.

```bash
go get github.com/permguard/sdk-go
```

It needs a Permguard data plane to ask: by default it listens on port `7443`, for HTTP and gRPC alike.
See the [Permguard repository](https://github.com/permguard/permguard) to run one.

## Quick start

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	permguard "github.com/permguard/sdk-go"
)

func main() {
	// grpc:// or grpcs:// speaks gRPC; http:// or https:// speaks JSON.
	client, err := permguard.NewClient("grpc://localhost:7443")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	response, err := client.Evaluate(context.Background(), &permguard.EvaluateRequest{
		Zone:      "acme",
		Ledger:    "main-ledger",
		Subject:   &permguard.Entity{Type: "User", ID: "alice"},
		Resource:  &permguard.Entity{Type: "Document", ID: "budget-2026"},
		Action:    &permguard.Action{Name: "read"},
		RequestID: "example-1",
	})
	if err != nil {
		var refusal *permguard.Refusal
		if errors.As(err, &refusal) {
			log.Fatalf("refused: class=%s code=%s message=%s", refusal.Class, refusal.Code, refusal.Message)
		}
		log.Fatal(err)
	}

	fmt.Println("permitted:", response.Decision)
}
```

The same program is in [examples/terminal](examples/terminal/cmd/main.go); set `PERMGUARD_PDP_URL` to point it at another data plane.

A request that names no `profile` is decided by the ledger's `default` profile: the manifest `permguard init` writes declares it, with one partition, `cedar`.
A ledger whose manifest declares other profiles is asked with the profile's name, and its partition inputs are keyed by that profile's partitions.

### The request

| Field             | What it says                                                            |
| ----------------- | ----------------------------------------------------------------------- |
| `Zone`, `Ledger`  | the ledger whose policies decide, in its zone                           |
| `Profile`         | the profile of the ledger's manifest that selects the partitions to run |
| `Subject`         | who asks: a type and an id, with optional properties                    |
| `Resource`        | what is acted on: a type and an id, with optional properties            |
| `Action`          | what is done: a name, with optional properties                          |
| `Context`         | any other fact the policies read                                        |
| `PartitionInputs` | runtime data for the partitions of the profile, see below               |
| `RequestID`       | your id for the call, given back in the answer                          |

### The answer

A deny is not an error: it is a successful response with `Decision == false`.
Validation, authorization, availability and server failures are returned as `*permguard.Refusal`, with their stable `Class` and `Code`.
A call that fails in the client is a `*permguard.Refusal` too, the same on HTTP and gRPC and in every Permguard SDK: class `unavailable` with code `timeout` when the call's timeout or the context's deadline passed, code `unreachable` when the data plane could not be reached.
A cancellation of your own context is not a refusal: the call returns the context's error.

## Partition inputs

Runtime data is addressed to the partition name the ledger's profile declares:

```go
request.PartitionInputs = permguard.PartitionInputs{
	"cedar": {
		Type: "permguard.cedar.entities.v1",
		Data: []any{
			map[string]any{
				"uid":     map[string]any{"type": "Team", "id": "engineering"},
				"attrs":   map[string]any{"active": true},
				"parents": []any{},
			},
		},
	},
}
```

The map key must be a partition name from the selected profile.
`Type` asserts the partition's input contract; it does not select a policy runtime.

## Many evaluations at once

Set the defaults at the top level, list the `Evaluations`, and call `EvaluateMany`:

```go
response, err := client.EvaluateMany(ctx, &permguard.EvaluateRequest{
	Zone:    "acme",
	Ledger:  "main-ledger",
	Subject: &permguard.Entity{Type: "User", ID: "alice"},
	Evaluations: []permguard.Evaluation{
		{
			Resource:  &permguard.Entity{Type: "Document", ID: "budget-2026"},
			Action:    &permguard.Action{Name: "read"},
			RequestID: "read",
		},
		{
			Resource:  &permguard.Entity{Type: "Document", ID: "budget-2026"},
			Action:    &permguard.Action{Name: "update"},
			RequestID: "update",
		},
	},
	Options: &permguard.EvaluationOptions{
		EvaluationsSemantic: permguard.ExecuteAll,
	},
})
```

Each evaluation inherits what it does not set from the request.

## Endpoints and options

| Endpoint   | Transport                                      |
| ---------- | ---------------------------------------------- |
| `http://`  | JSON over HTTP                                 |
| `https://` | JSON over HTTPS                                |
| `grpc://`  | `permguard.data.v1.PolicyDecisionPoint`        |
| `grpcs://` | `permguard.data.v1.PolicyDecisionPoint` on TLS |

An endpoint without a port uses the one its scheme implies, as the Permguard command line does: `443` over TLS, `80` in plaintext.
A data plane listens on `7443` by default, so name its port.

`NewClient` takes options after the endpoint:

| Option           | Sets                                                             |
| ---------------- | ---------------------------------------------------------------- |
| `WithTimeout`    | the timeout of every call, or the context's deadline when sooner |
| `WithTLSConfig`  | the TLS configuration, for `https://` and `grpcs://`             |
| `WithHTTPClient` | the `*http.Client` used for HTTP                                 |
| `WithHeader`     | a static header, also sent as gRPC metadata                      |

`client.GetConfiguration(ctx)` reads the data plane's discovery document.
The client reuses its HTTP transport or gRPC channel: create it once, and call `Close` when it is no longer needed; a second `Close` is no error.

## Credentials

In production a data plane admits a call only with an access token of its cluster's Authority Plane, issued to a service account whose scopes include `decision.evaluate` on that data plane.
Obtain one with the account's access key and secret key, then send it on every call as `authorization: Bearer <token>`, with the header option above.

```sh
curl -s -u "$PERMGUARD_ACCESS_KEY:$PERMGUARD_SECRET_KEY" \
  -d grant_type=client_credentials \
  https://permguard.example:8443/authority/v1/token | jq -r .access_token
```

A token lasts an hour unless asked shorter: obtain another before it expires.
The development configurations Permguard ships grant `decision.evaluate` to everybody, so the quick start needs no token.
See the [Authority Plane](https://docs.permguard.com/how-it-works/authority-plane) for service accounts, their scopes and keys.

## Compatibility

This major version implements `permguard.api.pdp.native.v1`.
Compatibility is tied to that versioned interface, not to a server minor version.

## Develop

```bash
make protoc
make test
make lint
```

`proto/v1/pdp.proto` is generated from the canonical contract in the Permguard server repository.
A change to the contract regenerates the Go bindings and passes both the HTTP and the gRPC contract tests.

## License

Apache License 2.0: see [LICENSE](LICENSE), [NOTICE.md](NOTICE.md) and [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

Permguard is a trademark of Nitro Agility S.r.l., see [TRADEMARKS.md](TRADEMARKS.md).
Go and the Go logo are trademarks of Google LLC, shown only to name the language this SDK is for.

# The official Go SDK for Permguard

[![GitHub License](https://img.shields.io/github/license/permguard/permguard-go)](https://github.com/permguard/permguard-go?tab=Apache-2.0-1-ov-file#readme)
[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/permguard)](https://x.com/intent/follow?original_referer=https%3A%2F%2Fdeveloper.x.com%2F&ref_src=twsrc%5Etfw%7Ctwcamp%5Ebuttonembed%7Ctwterm%5Efollow%7Ctwgr%5ETwitterDev&screen_name=Permguard)

[![Documentation](https://img.shields.io/website?label=Docs&url=https%3A%2F%2Fwww.permguard.com%2F)](https://www.permguard.com/)
[![Build, test and publish the artifacts](https://github.com/permguard/permguard-go/actions/workflows/permguard-go-ci.yml/badge.svg)](https://github.com/permguard/permguard-go/actions/workflows/permguard-go-ci.yml)

<p align="left">
  <img src="https://raw.githubusercontent.com/permguard/permguard-assets/main/pink-txt//1line.svg" class="center" width="400px" height="auto"/>
</p>

The Permguard GO SDK provides a simple and flexible client to perform authorization checks against a Permguard Policy Decision Point (PDP) service using gRPC.
Plase refer to the [Permguard Documentation](https://www.permguard.com/) for more information.

---

## Prerequisites

- **Go 1.23.5**

---

## Installation

Run the following command to install the SDK:

```bash
go get -u github.com/permguard/permguard-go
```

---

## Usage Example

Below is a sample Go code demonstrating how to create a Permguard client, build an authorization request using a builder pattern, and process the authorization response:

```go
import (
  "github.com/permguard/permguard-go"
  "github.com/permguard/permguard-go/az/azreq"
)

// Create a new Permguard client
azClient := permguard.NewAZClient(
  permguard.WithEndpoint("localhost", 9094),
)

// Create the Principal
principal := azreq.NewPrincipalBuilder("amy.smith@acmecorp.com").Build()

// Create the entities
entities := []map[string]any{
  {
      "uid": map[string]any{
      "type": "MagicFarmacia::Platform::BranchInfo",
      "id":   "subscription",
      },
      "attrs": map[string]any{
      "active": true,
    },
    "parents": []any{},
  },
}

// Create a new authorization request
req := azreq.NewAZAtomicRequestBuilder(273165098782, "fd1ac44e4afa4fc4beec622494d3175a",
  "amy.smith@acmecorp.com", "MagicFarmacia::Platform::Subscription", "MagicFarmacia::Platform::Action::create").
  // RequestID
  WithRequestID("1234").
  // Principal
  WithPrincipal(principal).
  // Entities
  WithEntitiesItems(azreq.CedarEntityKind, entities).
  // Subject
  WithSubjectKind(azreq.UserType).
  WithSubjectSource("keycloack").
  WithSubjectProperty("isSuperUser", true).
  // Resource
  WithResourceID("e3a786fd07e24bfa95ba4341d3695ae8").
  WithResourceProperty("isEnabled", true).
  // Action
  WithActionProperty("isEnabled", true).
  WithContextProperty("time", "2025-01-23T16:17:46+00:00").
  WithContextProperty("isSubscriptionActive", true).
  Build()

// Check the authorization
decsion, response, _ := azClient.Check(req)
if decsion {
  fmt.Println("✅ Authorization Permitted")
} else {
  fmt.Println("❌ Authorization Denied")
  if response.Context.ReasonAdmin != nil {
    fmt.Printf("-> Reason Admin: %s\n", response.Context.ReasonAdmin.Message)
  }
  if response.Context.ReasonUser != nil {
    fmt.Printf("-> Reason User: %s\n", response.Context.ReasonUser.Message)
  }
  for _, eval := range response.Evaluations {
    if eval.Context.ReasonUser != nil {
      fmt.Printf("-> Reason Admin: %s\n", eval.Context.ReasonAdmin.Message)
      fmt.Printf("-> Reason User: %s\n", eval.Context.ReasonUser.Message)
    }
  }
}
```

---

## Version Compatibility

Our SDK follows a versioning scheme aligned with the PermGuard server versions to ensure seamless integration. The versioning format is as follows:

**SDK Versioning Format:** `x.y.z`

- **x.y**: Indicates the compatible PermGuard server version.
- **z**: Represents the SDK's patch or minor updates specific to that server version.

**Compatibility Examples:**

- `SDK Version 1.3.0` is compatible with `PermGuard Server 1.3`.
- `SDK Version 1.3.1` includes minor improvements or bug fixes for `PermGuard Server 1.3`.

**Incompatibility Example:**

- `SDK Version 1.3.0` **may not be guaranteed** to be compatible with `PermGuard Server 1.4` due to potential changes introduced in server version `1.4`.

**Important:** Ensure that the major and minor versions (`x.y`) of the SDK match those of your PermGuard server to maintain compatibility.

---

Created by [Nitro Agility](https://www.nitroagility.com/).

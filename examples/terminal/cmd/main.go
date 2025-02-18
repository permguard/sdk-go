// Copyright 2024 Nitro Agility S.r.l.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/permguard/permguard-go"
)

func main() {
	// Create a new PermGuard client
	azClient := permguard.NewAZClient(
		permguard.WithPDPEndpoint("localhost", 9094),
	)

	// Create a new subject
	subject := permguard.NewSubjectBuilder("amy.smith@acmecorp.com").
		WithKind("user").
		WithSource("keycloack").
		WithProperty("isSuperUser", true).
		Build()

	// Create a new resource
	resource := permguard.NewResourceBuilder("MagicFarmacia::Platform::Subscription").
		WithID("e3a786fd07e24bfa95ba4341d3695ae8").
		WithProperty("isEnabled", true).
		Build()

	// Create a new action
	action := permguard.NewActionBuilder("MagicFarmacia::Platform::Action::view").
		WithProperty("isEnabled", true).
		Build()

	// Create a new Context
	context := permguard.NewContextBuilder().
		WithProperty("time", "2025-01-23T16:17:46+00:00").
		WithProperty("isSubscriptionActive", true).
		Build()

	// Create a new evaluation
	evaluation := permguard.NewAZEvaluationBuilder(subject, resource, action).
		WithRequestID("1234").
		WithContext(context).
		Build()

	// Create a new request
	req := permguard.NewAZRequestBuilder(273165098782, "fd1ac44e4afa4fc4beec622494d3175a").
		WithEvaluation(evaluation).
		Build()

	// Check the authorization
	decsion := azClient.Check(req)
	if decsion {
		fmt.Println("✅ Authorization Permitted")
	} else {
		fmt.Println("❌ Authorization Denied")
	}
}

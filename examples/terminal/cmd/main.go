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

func checkAtomicEvaluation() {
	// Create a new Permguard client
	azClient := permguard.NewAZClient(
		permguard.WithPDPEndpoint("localhost", 9094),
	)

	principal := permguard.NewPrincipalBuilder("amy.smith@acmecorp.com").Build()

	req := permguard.NewAZAtomicRequestBuilder(273165098782, "fd1ac44e4afa4fc4beec622494d3175a",
		"amy.smith@acmecorp.com", "MagicFarmacia::Platform::Subscription", "MagicFarmacia::Platform::Action::view").
		// RequestID
		WithRequestID("1234").
		// Principal
		WithPrincipal(principal).
		// Subject
		WithSubjectKind("user").
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
	decsion := azClient.Check(req)
	if decsion {
		fmt.Println("✅ Authorization Permitted")
	} else {
		fmt.Println("❌ Authorization Denied")
	}
}

// checkMultipleEvaluations checks multiple evaluations
func checkMultipleEvaluations() {
	// Create a new Permguard client
	azClient := permguard.NewAZClient(
		permguard.WithPDPEndpoint("localhost", 9094),
	)

	principal := permguard.NewPrincipalBuilder("amy.smith@acmecorp.com").Build()

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

	// Create ations
	actionView := permguard.NewActionBuilder("MagicFarmacia::Platform::Action::view").
		WithProperty("isEnabled", true).
		Build()
		
	actionCreate := permguard.NewActionBuilder("MagicFarmacia::Platform::Action::create").
		WithProperty("isEnabled", true).
		Build()

	// Create a new Context
	context := permguard.NewContextBuilder().
		WithProperty("time", "2025-01-23T16:17:46+00:00").
		WithProperty("isSubscriptionActive", true).
		Build()

	// Create evaluations
	evaluationView := permguard.NewAZEvaluationBuilder(subject, resource, actionView).
		WithRequestID("1234").
		WithPrincipal(principal).
		WithContext(context).
		Build()

	evaluationCreate := permguard.NewAZEvaluationBuilder(subject, resource, actionCreate).
		WithRequestID("7890").
		WithPrincipal(principal).
		WithContext(context).
		Build()

	// Create a new request
	req := permguard.NewAZRequestBuilder(273165098782, "fd1ac44e4afa4fc4beec622494d3175a").
		WithEvaluation(evaluationView).
		WithEvaluation(evaluationCreate).
		Build()

	// Check the authorization
	decsion := azClient.Check(req)
	if decsion {
		fmt.Println("✅ Authorization Permitted")
	} else {
		fmt.Println("❌ Authorization Denied")
	}
}

func main() {
	checkMultipleEvaluations()
}

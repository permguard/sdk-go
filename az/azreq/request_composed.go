// Copyright 2025 Nitro Agility S.r.l.
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

package azreq

const (
	// CedarEntityKind is the kind of the Cedar entity.
	CedarEntityKind = "cedar"
)

// AZRequest is the AZRequest object.
type AZRequest struct {
	zoneID          uint64
	policyStoreKind string
	policyStoreID   string
	principal       *Principal
	entitiesSchema  string
	entitiesItems   []map[string]any
	evaluations     []AZEvaluation
}

// GetZoneID returns the Zone ID of the AZRequest.
func (u *AZRequest) GetZoneID() uint64 {
	return u.zoneID
}

// GetPolicyStoreKind returns the policy store kind of the AZRequest.
func (u *AZRequest) GetPolicyStoreKind() string {
	return u.policyStoreKind
}

// GetPolicyStoreID returns the policy store ID of the AZRequest.
func (u *AZRequest) GetPolicyStoreID() string {
	return u.policyStoreID
}

// GetEntitiesSchema returns the entities schema of the AZRequest.
func (u *AZRequest) GetEntitiesSchema() string {
	return u.entitiesSchema
}

// GetPrincipal returns the principal of the AZEvaluation.
func (u *AZRequest) GetPrincipal() *Principal {
	return u.principal
}

// GetEntitiesItems returns the entities items of the AZRequest.
func (u *AZRequest) GetEntitiesItems() []map[string]any {
	return u.entitiesItems
}

// GetEvaluations returns the evaluations of the AZRequest.
func (u *AZRequest) GetEvaluations() []AZEvaluation {
	return u.evaluations
}

// AZRequestBuilder is the builder for the AZRequest object.
type AZRequestBuilder struct {
	azRequest *AZRequest
}

// NewAZRequestBuilder creates a new AZRequest builder.
func NewAZRequestBuilder(zoneID uint64, ledgerID string) *AZRequestBuilder {
	return &AZRequestBuilder{
		azRequest: &AZRequest{
			zoneID:          zoneID,
			policyStoreKind: "ledger",
			policyStoreID:   ledgerID,
			entitiesSchema:  "",
			entitiesItems:   []map[string]any{},
			evaluations:     []AZEvaluation{},
		},
	}
}

// WithPrincipal sets the principal of the AZEvaluation.
func (b *AZRequestBuilder) WithPrincipal(principal *Principal) *AZRequestBuilder {
	b.azRequest.principal = principal
	return b
}

// WithEntitiesMap sets the entities map to the AZRequest.
func (b *AZRequestBuilder) WithEntitiesMap(schema string, entities map[string]any) *AZRequestBuilder {
	b.azRequest.entitiesSchema = schema
	b.azRequest.entitiesItems = []map[string]any{entities}
	return b
}

// WithEntitiesItems sets the entities items to the AZRequest.
func (b *AZRequestBuilder) WithEntitiesItems(schema string, entities []map[string]any) *AZRequestBuilder {
	b.azRequest.entitiesSchema = schema
	b.azRequest.entitiesItems = entities
	if b.azRequest.entitiesItems == nil {
		b.azRequest.entitiesItems = []map[string]any{}
	}
	return b
}

// WithEvaluation adds an evaluation to the AZRequest.
func (b *AZRequestBuilder) WithEvaluation(evaluation *AZEvaluation) *AZRequestBuilder {
	b.azRequest.evaluations = append(b.azRequest.evaluations, *evaluation)
	return b
}

// Build builds the AZRequest object.
func (b *AZRequestBuilder) Build() *AZRequest {
	return b.azRequest
}

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
	// PolicyStoreKind is the kind of the Policy Store.
	PolicyStoreKind = "ledger"
	// CedarEntityKind is the kind of the Cedar entity.
	CedarEntityKind = "cedar"
)

// AZRequestBuilder is the builder for the AZRequest object.
type AZRequestBuilder struct {
	azRequest *AZRequest
}

// NewAZRequestBuilder creates a new AZRequest builder.
func NewAZRequestBuilder(zoneID uint64, ledgerID string) *AZRequestBuilder {
	return &AZRequestBuilder{
		azRequest: &AZRequest{
			AZModel: &AZModel{
				ZoneID: zoneID,
				PolicyStore: &PolicyStore{
					Kind: PolicyStoreKind,
					ID:   ledgerID,
				},
				Entities: &Entities{
					Schema: "",
					Items:  []map[string]any{},
				},
			},
			Evaluations: []Evaluation{},
		},
	}
}

// WithPrincipal sets the principal of the Evaluation.
func (b *AZRequestBuilder) WithPrincipal(principal *Principal) *AZRequestBuilder {
	b.azRequest.AZModel.Principal = principal
	return b
}

// WithEntitiesMap sets the entities map to the AZRequest.
func (b *AZRequestBuilder) WithEntitiesMap(schema string, entities map[string]any) *AZRequestBuilder {
	b.azRequest.AZModel.Entities.Schema = schema
	b.azRequest.AZModel.Entities.Items = []map[string]any{entities}
	return b
}

// WithEntitiesItems sets the entities items to the AZRequest.
func (b *AZRequestBuilder) WithEntitiesItems(schema string, entities []map[string]any) *AZRequestBuilder {
	b.azRequest.AZModel.Entities.Schema = schema
	b.azRequest.AZModel.Entities.Items = entities
	if b.azRequest.AZModel.Entities.Items == nil {
		b.azRequest.AZModel.Entities.Items = []map[string]any{}
	}
	return b
}

// WithEvaluation adds an evaluation to the AZRequest.
func (b *AZRequestBuilder) WithEvaluation(evaluation *Evaluation) *AZRequestBuilder {
	b.azRequest.Evaluations = append(b.azRequest.Evaluations, *evaluation)
	return b
}

// Build builds the AZRequest object.
func (b *AZRequestBuilder) Build() *AZRequest {
	instance := *b.azRequest
	return &instance
}

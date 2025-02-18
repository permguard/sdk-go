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

package permguard

// AZRequest is the AZRequest object.
type AZRequest struct {
	zoneID          uint64
	policyStoreType string
	policyStoreID   string
	evaluations     []AZEvaluation
}

// GetZoneID returns the Zone ID of the AZRequest.
func (u *AZRequest) GetZoneID() uint64 {
	return u.zoneID
}

// GetPolicyStoreType returns the policy store type of the AZRequest.
func (u *AZRequest) GetPolicyStoreType() string {
	return u.policyStoreType
}

// GetPolicyStoreID returns the policy store ID of the AZRequest.
func (u *AZRequest) GetPolicyStoreID() string {
	return u.policyStoreID
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
			policyStoreType: "ledger",
			policyStoreID:   ledgerID,
			evaluations:     []AZEvaluation{},
		},
	}
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

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

package azreq

const (
	UserType      = "user"
	WorkloadType  = "workload"
	AttributeType = "attribute"
)

// PolicyStore represents the Policy Store.
type PolicyStore struct {
	Kind string `json:"kind,omitempty"`
	ID   string `json:"id,omitempty"`
}

// Entities represents the Entities.
type Entities struct {
	Schema string           `json:"schema,omitempty"`
	Items  []map[string]any `json:"items,omitempty"`
}

// Evaluation is the Evaluation object.
type Evaluation struct {
	RequestID string         `json:"request_id,omitempty"`
	Subject   *Subject       `json:"subject,omitempty"`
	Resource  *Resource      `json:"resource,omitempty"`
	Action    *Action        `json:"action,omitempty"`
	Context   map[string]any `json:"context,omitempty"`
}

// AZModel is the Authorization Model.
type AZModel struct {
	ZoneID      uint64       `json:"zone_id"`
	Principal   *Principal   `json:"principal,omitempty"`
	PolicyStore *PolicyStore `json:"policy_store,omitempty"`
	Entities    *Entities    `json:"entities,omitempty"`
}

// AZRequest is the AZRequest object.
type AZRequest struct {
	AZModel     *AZModel  `json:"authorization_model,omitempty"`
	RequestID   string    `json:"request_id,omitempty"`
	Subject     *Subject  `json:"subject,omitempty"`
	Resource    *Resource `json:"resource,omitempty"`
	Action      *Action   `json:"action,omitempty"`
	Context     map[string]any
	Evaluations []Evaluation `json:"evaluations,omitempty"`
}

// Principal is the principal object.
type Principal struct {
	Type   string `json:"type,omitempty"`
	ID     string `json:"id,omitempty"`
	Source string `json:"source,omitempty"`
}

// Subject is the subject object.
type Subject struct {
	Type       string         `json:"type,omitempty"`
	ID         string         `json:"id,omitempty"`
	Source     string         `json:"source,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

// Resource is the resource object.
type Resource struct {
	Type       string         `json:"type,omitempty"`
	ID         string         `json:"id,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

// Action is the action object.
type Action struct {
	Name       string         `json:"name,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

// ReasonResponse provides the rationale for the response.
type ReasonResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ContextResponse represents the context included in the response.
type ContextResponse struct {
	ID          string          `json:"id,omitempty"`
	ReasonAdmin *ReasonResponse `json:"reason_admin,omitempty"`
	ReasonUser  *ReasonResponse `json:"reason_user,omitempty"`
}

// EvaluationResponse represents the result of the evaluation process.
type EvaluationResponse struct {
	RequestID string           `json:"request_id"`
	Decision  bool             `json:"decision,omitempty"`
	Context   *ContextResponse `json:"context,omitempty"`
}

// AZResponse represents the outcome of the authorization decision.
type AZResponse struct {
	RequestID   string               `json:"request_id,omitempty"`
	Decision    bool                 `json:"decision"`
	Context     *ContextResponse     `json:"context,omitempty"`
	Evaluations []EvaluationResponse `json:"evaluations,omitempty"`
}

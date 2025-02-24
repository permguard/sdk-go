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

// PolicyStore represents the Policy Store.
type PolicyStore struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

// Entities represents the Entities.
type Entities struct {
	Schema string           `json:"schema"`
	Items  []map[string]any `json:"items"`
}

// Evaluation is the Evaluation object.
type Evaluation struct {
	RequestID string         `json:"request_id"`
	Subject   *Subject       `json:"subject"`
	Resource  *Resource      `json:"resource"`
	Action    *Action        `json:"action"`
	Context   map[string]any `json:"context"`
}

// AZModel is the Authorization Model.
type AZModel struct {
	ZoneID      uint64       `json:"zone_id"`
	Principal   *Principal   `json:"principal"`
	PolicyStore *PolicyStore `json:"policy_store"`
	Entities    *Entities    `json:"entities"`
}

// AZRequest is the AZRequest object.
type AZRequest struct {
	AZModel     *AZModel     `json:"authorization_model"`
	Evaluations []Evaluation `json:"evaluations"`
}

// Principal is the principal object.
type Principal struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Source string `json:"source"`
}

// Subject is the subject object.
type Subject struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Source     string         `json:"source"`
	Properties map[string]any `json:"properties"`
}

// Resource is the resource object.
type Resource struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Properties map[string]any `json:"properties"`
}

// Action is the action object.
type Action struct {
	Name       string         `json:"name"`
	Properties map[string]any `json:"properties"`
}

// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

// EvaluationsSemantic controls how a boxcarred request is executed and combined.
type EvaluationsSemantic string

const (
	ExecuteAll          EvaluationsSemantic = "execute_all"
	DenyOnFirstDeny     EvaluationsSemantic = "deny_on_first_deny"
	PermitOnFirstPermit EvaluationsSemantic = "permit_on_first_permit"
)

// Entity identifies a subject, resource, or caller.
type Entity struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Properties map[string]any `json:"properties,omitempty"`
}

// Action is the operation being evaluated.
type Action struct {
	Name       string         `json:"name"`
	Properties map[string]any `json:"properties,omitempty"`
}

// PartitionInput supplies runtime data to one named profile partition.
type PartitionInput struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

// PartitionInputs are addressed by the partition name declared by the profile.
type PartitionInputs map[string]PartitionInput

// Evaluation overrides the request defaults for one item in a boxcarred request.
type Evaluation struct {
	Subject         *Entity          `json:"subject,omitempty"`
	Resource        *Entity          `json:"resource,omitempty"`
	Action          *Action          `json:"action,omitempty"`
	Context         *map[string]any  `json:"context,omitempty"`
	PartitionInputs *PartitionInputs `json:"partition_inputs,omitempty"`
	RequestID       string           `json:"request_id,omitempty"`
}

// EvaluationOptions controls boxcarred evaluation behavior.
type EvaluationOptions struct {
	EvaluationsSemantic EvaluationsSemantic `json:"evaluations_semantic,omitempty"`
}

// EvaluateRequest is the payload of permguard.api.pdp.native.v1.
type EvaluateRequest struct {
	Zone            string             `json:"zone"`
	Ledger          string             `json:"ledger"`
	Profile         string             `json:"profile,omitempty"`
	Subject         *Entity            `json:"subject,omitempty"`
	Resource        *Entity            `json:"resource,omitempty"`
	Action          *Action            `json:"action,omitempty"`
	Context         map[string]any     `json:"context,omitempty"`
	Principal       *Entity            `json:"principal,omitempty"`
	PartitionInputs PartitionInputs    `json:"partition_inputs,omitempty"`
	Evaluations     []Evaluation       `json:"evaluations,omitempty"`
	Options         *EvaluationOptions `json:"options,omitempty"`
	RequestID       string             `json:"request_id,omitempty"`
}

// Reason is a stable machine code and its human-readable explanation.
type Reason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// DecisionContext carries the evidence associated with one decision.
type DecisionContext struct {
	ID           string   `json:"id,omitempty"`
	ReasonAdmin  *Reason  `json:"reason_admin,omitempty"`
	ReasonUser   *Reason  `json:"reason_user,omitempty"`
	Policies     []string `json:"policies,omitempty"`
	AbsentInputs []string `json:"absent_inputs,omitempty"`
}

// Decision is one item in a boxcarred response.
type Decision struct {
	Decision  bool             `json:"decision"`
	RequestID string           `json:"request_id,omitempty"`
	Context   *DecisionContext `json:"context,omitempty"`
}

// EvaluateResponse is the PDP's decision response.
type EvaluateResponse struct {
	Decision    bool             `json:"decision"`
	RequestID   string           `json:"request_id,omitempty"`
	Context     *DecisionContext `json:"context,omitempty"`
	Evaluations []Decision       `json:"evaluations,omitempty"`
}

// Endpoints are the HTTP bindings advertised by the PDP.
type Endpoints struct {
	Evaluation  string `json:"evaluation"`
	Evaluations string `json:"evaluations"`
}

// StoreScope describes where and how a request names its policy store.
type StoreScope struct {
	In      string `json:"in"`
	Zone    string `json:"zone"`
	Ledger  string `json:"ledger"`
	Profile string `json:"profile"`
}

// Configuration describes the PDP native v1 interface served by an endpoint.
type Configuration struct {
	Interface    string     `json:"interface"`
	PDP          string     `json:"pdp"`
	Endpoints    Endpoints  `json:"endpoints"`
	Capabilities []string   `json:"capabilities"`
	StoreScope   StoreScope `json:"store_scope"`
}

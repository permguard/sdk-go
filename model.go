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
	// RequestNonce asks for a signed response: base64url without padding, 16 to 64 bytes.
	// Empty asks for none; a plane that cannot sign refuses rather than answer unsigned.
	RequestNonce string `json:"request_nonce,omitempty"`
}

// Reason is a stable machine code and its human-readable explanation.
type Reason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// EvidencePosition is where a decision record sits in its plane's stream.
type EvidencePosition struct {
	HostID   string `json:"host_id"`
	Run      string `json:"run"`
	Sequence string `json:"sequence"`
}

// Evidence says whether and where the decision was recorded.
type Evidence struct {
	Mode     string            `json:"mode"`
	Status   string            `json:"status"`
	RecordID *string           `json:"record_id"`
	Position *EvidencePosition `json:"position"`
}

// DecisionContext carries the evidence associated with one decision.
type DecisionContext struct {
	ID           string    `json:"id,omitempty"`
	ReasonAdmin  *Reason   `json:"reason_admin,omitempty"`
	ReasonUser   *Reason   `json:"reason_user,omitempty"`
	Policies     []string  `json:"policies,omitempty"`
	AbsentInputs []string  `json:"absent_inputs,omitempty"`
	Evidence     *Evidence `json:"evidence,omitempty"`
}

// Attestation is the plane's signature over a response it was asked to sign: a flattened JWS.
type Attestation struct {
	Protected string `json:"protected"`
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
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
	Attestation *Attestation     `json:"attestation,omitempty"`
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

// ResponseSigning says whether the PDP signs a response a request asks it to sign.
type ResponseSigning struct {
	State            string `json:"state"`
	ValiditySeconds  uint32 `json:"validity_seconds,omitempty"`
	ClockSkewSeconds uint32 `json:"clock_skew_seconds,omitempty"`
}

// Configuration describes the PDP native v1 interface served by an endpoint.
type Configuration struct {
	Interface       string            `json:"interface"`
	Version         string            `json:"version"`
	Component       string            `json:"component"`
	State           string            `json:"state"`
	Stage           string            `json:"stage"`
	Sunset          *string           `json:"sunset,omitempty"`
	PDP             string            `json:"pdp"`
	Endpoints       Endpoints         `json:"endpoints"`
	Capabilities    []string          `json:"capabilities"`
	StoreScope      StoreScope        `json:"store_scope"`
	Limits          map[string]uint64 `json:"limits"`
	ResponseSigning ResponseSigning   `json:"response_signing"`
}

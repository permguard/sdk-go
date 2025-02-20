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

// AZEvaluation is the AZEvaluation object.
type AZEvaluation struct {
	requestID string
	subject   *Subject
	resource  *Resource
	action    *Action
	context   map[string]any
}

// GetRequestID returns the request ID of the AZEvaluation.
func (u *AZEvaluation) GetRequestID() string {
	return u.requestID
}

// GetSubject returns the subject of the AZEvaluation.
func (u *AZEvaluation) GetSubject() *Subject {
	return u.subject
}

// GetResource returns the resource of the AZEvaluation.
func (u *AZEvaluation) GetResource() *Resource {
	return u.resource
}

// GetAction returns the action of the AZEvaluation.
func (u *AZEvaluation) GetAction() *Action {
	return u.action
}

// GetContext returns the context of the AZEvaluation.
func (u *AZEvaluation) GetContext() map[string]any {
	return u.context
}

// AZEvaluationBuilder is the builder for the AZEvaluation object.
type AZEvaluationBuilder struct {
	azEvaluation *AZEvaluation
}

// NewAZEvaluationBuilder creates a new AZEvaluation builder.
func NewAZEvaluationBuilder(subject *Subject, resource *Resource, action *Action) *AZEvaluationBuilder {
	return &AZEvaluationBuilder{
		azEvaluation: &AZEvaluation{
			subject:  subject,
			resource: resource,
			action:   action,
		},
	}
}

// WithRequestID sets the request ID of the AZEvaluation.
func (b *AZEvaluationBuilder) WithRequestID(requestID string) *AZEvaluationBuilder {
	b.azEvaluation.requestID = requestID
	return b
}

// WithContext sets the context of the AZEvaluation.
func (b *AZEvaluationBuilder) WithContext(context map[string]any) *AZEvaluationBuilder {
	b.azEvaluation.context = context
	return b
}

// Build builds the AZEvaluation object.
func (b *AZEvaluationBuilder) Build() *AZEvaluation {
	return b.azEvaluation
}

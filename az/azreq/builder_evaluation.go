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

// EvaluationBuilder is the builder for the Evaluation object.
type EvaluationBuilder struct {
	azEvaluation *Evaluation
}

// NewEvaluationBuilder creates a new Evaluation builder.
func NewEvaluationBuilder(subject *Subject, resource *Resource, action *Action) *EvaluationBuilder {
	return &EvaluationBuilder{
		azEvaluation: &Evaluation{
			Subject:  subject,
			Resource: resource,
			Action:   action,
		},
	}
}

// WithRequestID sets the request ID of the Evaluation.
func (b *EvaluationBuilder) WithRequestID(requestID string) *EvaluationBuilder {
	b.azEvaluation.RequestID = requestID
	return b
}

// WithContext sets the context of the Evaluation.
func (b *EvaluationBuilder) WithContext(context map[string]any) *EvaluationBuilder {
	b.azEvaluation.Context = context
	return b
}

// Build builds the Evaluation object.
func (b *EvaluationBuilder) Build() *Evaluation {
	instance := *b.azEvaluation
	instance.Context = deepCopyMap(instance.Context)
	return &instance
}

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

// AZAtomicRequestBuilder is the builder for the AZAtomicRequest object.
type AZAtomicRequestBuilder struct {
	requestID         string
	azSubjectBuilder  *SubjectBuilder
	azResourceBuilder *ResourceBuilder
	azActionBuilder   *ActionBuilder
	azContextBuilder  *ContextBuilder
	azRequestBuilder  *AZRequestBuilder
}

// NewAZAtomicRequestBuilder creates a new AZAtomicRequest builder.
func NewAZAtomicRequestBuilder(zoneID uint64, ledgerID, subjectID, resourceKind, actionName string) *AZAtomicRequestBuilder {
	azReqBuilder := NewAZRequestBuilder(zoneID, ledgerID)
	azSubjectBuilder := NewSubjectBuilder(subjectID)
	azResourceBuilder := NewResourceBuilder(resourceKind)
	azActionBuilder := NewActionBuilder(actionName)
	azContextBuilder := NewContextBuilder()
	return &AZAtomicRequestBuilder{
		azRequestBuilder:  azReqBuilder,
		azSubjectBuilder:  azSubjectBuilder,
		azResourceBuilder: azResourceBuilder,
		azActionBuilder:   azActionBuilder,
		azContextBuilder:  azContextBuilder,
	}
}

// WithRequestID sets the ID of the request.
func (b *AZAtomicRequestBuilder) WithRequestID(requestID string) *AZAtomicRequestBuilder {
	b.requestID = requestID
	return b
}

// WithSubjectSource sets the source of the subject.
func (b *AZAtomicRequestBuilder) WithSubjectKind(kind string) *AZAtomicRequestBuilder {
	b.azSubjectBuilder.WithKind(kind)
	return b
}

// WithSubjectSource sets the source of the subject.
func (b *AZAtomicRequestBuilder) WithSubjectSource(source string) *AZAtomicRequestBuilder {
	b.azSubjectBuilder.WithSource(source)
	return b
}

// WithSubjectProperty sets a property of the subject.
func (b *AZAtomicRequestBuilder) WithSubjectProperty(key string, value interface{}) *AZAtomicRequestBuilder {
	b.azSubjectBuilder.WithProperty(key, value)
	return b
}

// WithResourceID sets the ID of the resource.
func (b *AZAtomicRequestBuilder) WithResourceID(id string) *AZAtomicRequestBuilder {
	b.azResourceBuilder.WithID(id)
	return b
}

// WithResourceProperty sets a property of the resource.
func (b *AZAtomicRequestBuilder) WithResourceProperty(key string, value interface{}) *AZAtomicRequestBuilder {
	b.azResourceBuilder.WithProperty(key, value)
	return b
}

// WithActionProperty sets a property of the action.
func (b *AZAtomicRequestBuilder) WithActionProperty(key string, value interface{}) *AZAtomicRequestBuilder {
	b.azActionBuilder.WithProperty(key, value)
	return b
}

// WithContextProperty sets a property of the context.
func (b *AZAtomicRequestBuilder) WithContextProperty(key string, value interface{}) *AZAtomicRequestBuilder {
	b.azContextBuilder.WithProperty(key, value)
	return b
}

// Build builds the AZAtomicRequest object.
func (b *AZAtomicRequestBuilder) Build() *AZRequest {
	subject := b.azSubjectBuilder.Build()
	resource := b.azResourceBuilder.Build()
	action := b.azActionBuilder.Build()
	context := b.azContextBuilder.Build()
	azEvaluation := NewAZEvaluationBuilder(subject, resource, action).
		WithContext(context).
		Build()
	b.azRequestBuilder.WithEvaluation(azEvaluation)
	return b.azRequestBuilder.Build()
}

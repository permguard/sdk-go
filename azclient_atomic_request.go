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
	principal         *Principal
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

// WithEntitiesMap sets the entities map to the AZRequest.
func (b *AZAtomicRequestBuilder) WithEntitiesMap(schema string, entities map[string]any) *AZAtomicRequestBuilder {
	b.WithEntitiesMap(schema, entities)
	return b
}

// WithEntitiesItems sets the entities items to the AZRequest.
func (b *AZAtomicRequestBuilder) WithEntitiesItems(schema string, entities []map[string]any) *AZAtomicRequestBuilder {
	b.WithEntitiesItems(schema, entities)
	return b
}

// WithRequestID sets the ID of the AZRequest.
func (b *AZAtomicRequestBuilder) WithRequestID(requestID string) *AZAtomicRequestBuilder {
	b.requestID = requestID
	return b
}

// WithPrincipal sets the principal of the AZRequest.
func (b *AZAtomicRequestBuilder) WithPrincipal(principal *Principal) *AZAtomicRequestBuilder {
	b.principal = principal
	return b
}

// WithSubjectKind sets the kind of the subject for the AZRequest.
func (b *AZAtomicRequestBuilder) WithSubjectKind(kind string) *AZAtomicRequestBuilder {
	b.azSubjectBuilder.WithKind(kind)
	return b
}

// WithSubjectSource sets the source of the subject for the AZRequest.
func (b *AZAtomicRequestBuilder) WithSubjectSource(source string) *AZAtomicRequestBuilder {
	b.azSubjectBuilder.WithSource(source)
	return b
}

// WithSubjectProperty sets a property of the subject for the AZRequest.
func (b *AZAtomicRequestBuilder) WithSubjectProperty(key string, value interface{}) *AZAtomicRequestBuilder {
	b.azSubjectBuilder.WithProperty(key, value)
	return b
}

// WithResourceID sets the ID of the resource for the AZRequest.
func (b *AZAtomicRequestBuilder) WithResourceID(id string) *AZAtomicRequestBuilder {
	b.azResourceBuilder.WithID(id)
	return b
}

// WithResourceProperty sets a property of the resource for the AZRequest.
func (b *AZAtomicRequestBuilder) WithResourceProperty(key string, value interface{}) *AZAtomicRequestBuilder {
	b.azResourceBuilder.WithProperty(key, value)
	return b
}

// WithActionProperty sets a property of the action for the AZRequest.
func (b *AZAtomicRequestBuilder) WithActionProperty(key string, value interface{}) *AZAtomicRequestBuilder {
	b.azActionBuilder.WithProperty(key, value)
	return b
}

// WithContextProperty sets a property of the context for the AZRequest.
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

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

// PrincipalBuilder is the builder for the principal object.
type PrincipalBuilder struct {
	principal *Principal
}

// NewPrincipalBuilder creates a new principal builder.
func NewPrincipalBuilder(id string) *PrincipalBuilder {
	return &PrincipalBuilder{
		principal: &Principal{
			ID:   id,
			Type: UserType,
		},
	}
}

// WithUserType sets as user the type of the subject for the AZRequest.
func (b *PrincipalBuilder) WithUserType() *PrincipalBuilder {
	b.principal.Type = UserType
	return b
}

// WithWorkloadType sets as workload the type of the subject for the AZRequest.
func (b *PrincipalBuilder) WithWorkloadType() *PrincipalBuilder {
	b.principal.Type = WorkloadType
	return b
}

// WithAttributeType sets as attribute the type of the subject for the AZRequest.
func (b *PrincipalBuilder) WithAttributeType() *PrincipalBuilder {
	b.principal.Type = AttributeType
	return b
}

// WithType sets the type of the principal.
func (b *PrincipalBuilder) WithType(principalType string) *PrincipalBuilder {
	b.principal.Type = principalType
	return b
}

// WithProperty sets the property of the principal.
func (b *PrincipalBuilder) WithSource(source string) *PrincipalBuilder {
	b.principal.Source = source
	return b
}

// Build builds the principal object.
func (b *PrincipalBuilder) Build() *Principal {
	instance := *b.principal
	return &instance
}

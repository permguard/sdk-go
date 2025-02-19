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

// Resource is the resource object.
type Resource struct {
	kind       string
	id         string
	properties map[string]any
}

// GetKind returns the kind of the resource.
func (u *Resource) GetKind() string {
	return u.kind
}

// GetID returns the ID of the resource.
func (u *Resource) GetID() string {
	return u.id
}

// GetProperties returns the properties of the resource.
func (u *Resource) GetProperties() map[string]any {
	return u.properties
}

// ResourceBuilder is the builder for the resource object.
type ResourceBuilder struct {
	resource *Resource
}

// NewResourceBuilder creates a new resource builder.
func NewResourceBuilder(kind string) *ResourceBuilder {
	return &ResourceBuilder{
		resource: &Resource{
			kind: kind,
		},
	}
}

// WithID sets the id of the resource.
func (b *ResourceBuilder) WithID(kind string) *ResourceBuilder {
	b.resource.kind = kind
	return b
}

// WithProperty sets the property of the resource.
func (b *ResourceBuilder) WithProperty(key string, value any) *ResourceBuilder {
	if b.resource.properties == nil {
		b.resource.properties = make(map[string]any)
	}
	b.resource.properties[key] = value
	return b
}

// Build builds the resource object.
func (b *ResourceBuilder) Build() *Resource {
	return b.resource
}

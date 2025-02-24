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

// ResourceBuilder is the builder for the resource object.
type ResourceBuilder struct {
	resource *Resource
}

// NewResourceBuilder creates a new resource builder.
func NewResourceBuilder(kind string) *ResourceBuilder {
	return &ResourceBuilder{
		resource: &Resource{
			Type: kind,
		},
	}
}

// WithID sets the id of the resource.
func (b *ResourceBuilder) WithID(id string) *ResourceBuilder {
	b.resource.ID = id
	return b
}

// WithProperty sets the property of the resource.
func (b *ResourceBuilder) WithProperty(key string, value any) *ResourceBuilder {
	if b.resource.Properties == nil {
		b.resource.Properties = make(map[string]any)
	}
	b.resource.Properties[key] = value
	return b
}

// Build builds the resource object.
func (b *ResourceBuilder) Build() *Resource {
	instance := *b.resource
	instance.Properties = deepCopyMap(instance.Properties)
	return &instance
}

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

// ContextBuilder is the builder for the context object.
type ContextBuilder struct {
	context map[string]any
}

// NewContextBuilder creates a new context builder.
func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{
		context: map[string]any{},
	}
}

// WithProperty sets the property of the context.
func (b *ContextBuilder) WithProperty(key string, value any) *ContextBuilder {
	if b.context == nil {
		b.context = map[string]any{}
	}
	b.context[key] = value
	return b
}

// Build builds the context object.
func (b *ContextBuilder) Build() map[string]any {
	instance := deepCopyMap(b.context)
	return instance
}

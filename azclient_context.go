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

// Context is the context object.
type Context struct {
	properties map[string]any
}

// GetProperties returns the properties of the context.
func (u *Context) GetProperties() map[string]any {
	return u.properties
}

// ContextBuilder is the builder for the context object.
type ContextBuilder struct {
	context *Context
}

// NewContextBuilder creates a new context builder.
func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{
		context: &Context{},
	}
}

// WithProperty sets the property of the context.
func (b *ContextBuilder) WithProperty(key string, value any) *ContextBuilder {
	if b.context.properties == nil {
		b.context.properties = make(map[string]any)
	}
	b.context.properties[key] = value
	return b
}

// Build builds the context object.
func (b *ContextBuilder) Build() *Context {
	return b.context
}

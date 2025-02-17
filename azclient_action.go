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

// Action is the action object.
type Action struct {
	name       string
	properties map[string]any
}

// GetName returns the name of the action.
func (u *Action) GetName() string {
	return u.name
}

// GetProperties returns the properties of the action.
func (u *Action) GetProperties() map[string]any {
	return u.properties
}

// ActionBuilder is the builder for the action object.
type ActionBuilder struct {
	action *Action
}

// NewActionBuilder creates a new action builder.
func NewActionBuilder(name string) *ActionBuilder {
	return &ActionBuilder{
		action: &Action{
			name: name,
		},
	}
}

// WithProperty sets the property of the action.
func (b *ActionBuilder) WithProperty(key string, value any) *ActionBuilder {
	if b.action.properties == nil {
		b.action.properties = make(map[string]any)
	}
	b.action.properties[key] = value
	return b
}

// Build builds the action object.
func (b *ActionBuilder) Build() *Action {
	return b.action
}

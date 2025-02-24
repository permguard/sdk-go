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

// Evaluation is the Evaluation object.
type Evaluation struct {
	requestID string
	subject   *Subject
	resource  *Resource
	action    *Action
	context   map[string]any
}

// GetRequestID returns the request ID of the Evaluation.
func (u *Evaluation) GetRequestID() string {
	return u.requestID
}

// GetSubject returns the subject of the Evaluation.
func (u *Evaluation) GetSubject() *Subject {
	return u.subject
}

// GetResource returns the resource of the Evaluation.
func (u *Evaluation) GetResource() *Resource {
	return u.resource
}

// GetAction returns the action of the Evaluation.
func (u *Evaluation) GetAction() *Action {
	return u.action
}

// GetContext returns the context of the Evaluation.
func (u *Evaluation) GetContext() map[string]any {
	return u.context
}

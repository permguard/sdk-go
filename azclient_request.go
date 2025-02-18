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

// AZRequest is the AZRequest object.
type AZRequest struct {
	requestID string
	subject   *Subject
	resource  *Resource
	action    *Action
}

// GetRequestID returns the request ID of the AZRequest.
func (u *AZRequest) GetRequestID() string {
	return u.requestID
}

// GetSubject returns the subject of the AZRequest.
func (u *AZRequest) GetSubject() *Subject {
	return u.subject
}

// GetResource returns the resource of the AZRequest.
func (u *AZRequest) GetResource() *Resource {
	return u.resource
}

// GetAction returns the action of the AZRequest.
func (u *AZRequest) GetAction() *Action {
	return u.action
}

// AZRequestBuilder is the builder for the AZRequest object.
type AZRequestBuilder struct {
	AZRequest *AZRequest
}

// NewAZRequestBuilder creates a new AZRequest builder.
func NewAZRequestBuilder(subject *Subject, resource *Resource, action *Action) *AZRequestBuilder {
	return &AZRequestBuilder{
		AZRequest: &AZRequest{
			subject:  subject,
			resource: resource,
			action:   action,
		},
	}
}

// WithRequestID sets the request ID of the AZRequest.
func (b *AZRequestBuilder) WithRequestID(requestID string) *AZRequestBuilder {
	b.AZRequest.requestID = requestID
	return b
}

// Build builds the AZRequest object.
func (b *AZRequestBuilder) Build() *AZRequest {
	return b.AZRequest
}

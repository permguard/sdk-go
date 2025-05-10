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

import (
	"fmt"

	"github.com/permguard/sdk-go/az/azreq"
	v1 "github.com/permguard/sdk-go/internal/az/azreq/grpc/v1"
)

// AZClient is the client to interact with the authorization server.
type AZClient struct {
	azConfig *AZConfig
}

// NewAZClient creates a new authorization client.
func NewAZClient(opts ...AZOption) *AZClient {
	c := &AZClient{}
	c.azConfig = &AZConfig{
		pdpEndpoint: &azEndpoint{
			endpoint: "localhost",
			port:     9094,
		},
	}
	for _, opt := range opts {
		opt(c.azConfig)
	}
	return c
}

// Check checks the input authorization request with the authorization server.
func (c *AZClient) Check(req *azreq.AZRequest) (bool, *azreq.AZResponse, error) {
	target := fmt.Sprintf("%s:%d", c.azConfig.pdpEndpoint.endpoint, c.azConfig.pdpEndpoint.port)
	canExecute, err := v1.AuthorizationCheck(target, req)
	decision := false
	if canExecute != nil {
		decision = canExecute.Decision
	}
	return decision, canExecute, err
}

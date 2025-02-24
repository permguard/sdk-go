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
	"context"
	"fmt"

	"github.com/permguard/permguard-go/az/azreq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/permguard/permguard-go/internal/az/azreq/grpc/v1"
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
func (c *AZClient) Check(req *azreq.AZRequest) bool {
	target := fmt.Sprintf("%s:%d", c.azConfig.pdpEndpoint.endpoint, c.azConfig.pdpEndpoint.port)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false
	}
	client := v1.NewV1PDPServiceClient(conn)
	azCheckRequest := &v1.AuthorizationCheckRequest{}
	
	azCheckRequest.AuthorizationModel = &v1.AuthorizationModelRequest{}
	azCheckRequest.AuthorizationModel.ZoneID = int64(req.GetAuthZModel().GetZoneID())
	azCheckRequest.AuthorizationModel.PolicyStore = &v1.PolicyStore{}
	azCheckRequest.AuthorizationModel.PolicyStore.ID = "abc"
	azCheckRequest.AuthorizationModel.PolicyStore.Kind = "ledger"

	azCheckRequest.Subject = &v1.Subject{}
	azCheckRequest.Subject.ID = "pippo"
	azCheckRequest.Subject.Kind = "user"

	azCheckRequest.Resource = &v1.Resource{}
	azCheckRequest.Resource.ID = "pippo"
	azCheckRequest.Resource.Kind = "user"

	azCheckRequest.Action = &v1.Action{}
	azCheckRequest.Action.Name = "pippo"

	ctx := context.Background()
	azCheckResponse, err := client.AuthorizationCheck(ctx, azCheckRequest)
	if err != nil {
		return false
	}
	return azCheckResponse.Decision
}

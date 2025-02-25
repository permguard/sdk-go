// Copyright 2024 Nitro Agility S.r.l.
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

package v1

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/permguard/permguard-go/az/azreq"
)

// AuthorizationCheck checks the authorization request with the authorization server.
func AuthorizationCheck(endpoint string, req *azreq.AZRequest) (*azreq.AZResponse, error) {
	if req == nil || req.Evaluations == nil {
		return nil, errors.New("pep: invalid request")
	}
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, fmt.Errorf("pep: grpc connection error: %w", err)
	}

	client := NewV1PDPServiceClient(conn)
	azCheckRequest, err := MapAZRequestToGrpcAuthorizationCheckRequest(req)
	if err != nil {
		return nil, fmt.Errorf("pep: grpc request could not be constructed error: %w", err)
	}

	ctx := context.Background()
	azCheckResponse, err := client.AuthorizationCheck(ctx, azCheckRequest)
	if err != nil {
		return nil, fmt.Errorf("pep: grpc call failed error: %w", err)
	}
	azResponse, err := MapGrpcAuthorizationCheckResponseToAZResponse(azCheckResponse)
	return azResponse, err
}

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

// azEndpoint is the endpoint for the authorization server.
type azEndpoint struct {
	endpoint string
	port     int
}

// AZOption is the option for the authorization client configuration.
type AZOption func(*AZConfig)

// AZConfig is the configuration for the authorization client.
type AZConfig struct {
	pdpEndpoint *azEndpoint
}

// WithEndpoint sets the gRPC endpoint for the authorization server.
func WithEndpoint(endpoint string, port int) AZOption {
	return func(c *AZConfig) {
		c.pdpEndpoint = &azEndpoint{
			endpoint: endpoint,
			port:     port,
		}
	}
}

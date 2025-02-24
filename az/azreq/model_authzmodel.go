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

// AuthZModel is the Authorization Model.
type AuthZModel struct {
	zoneID      uint64
	principal   *Principal
	policyStore *PolicyStore
	entities    *Entities
}

// AZRequest is the AZRequest object.
func (u *AuthZModel) GetZoneID() uint64 {
	return u.zoneID
}

// GetPrincipal returns the principal of the AuthZModel.
func (u *AuthZModel) GetPrincipal() *Principal {
	return u.principal
}

// GetPolicyStore returns the policy store of the AuthZModel.
func (u *AuthZModel) GetPolicyStore() *PolicyStore {
	return u.policyStore
}

// GetEntities returns the entities of the AuthZModel.
func (u *AuthZModel) GetEntities() *Entities {
	return u.entities
}

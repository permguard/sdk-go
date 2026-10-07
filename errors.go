// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

// Refusal is a structured PDP error. A deny is not a Refusal: it is a
// successful EvaluateResponse whose Decision is false.
type Refusal struct {
	Class      string     `json:"class"`
	Code       string     `json:"code"`
	Message    string     `json:"message"`
	HTTPStatus int        `json:"-"`
	GRPCCode   codes.Code `json:"-"`
}

func (r *Refusal) Error() string {
	if r == nil {
		return ""
	}
	if r.Code == "" {
		return r.Message
	}
	return fmt.Sprintf("%s: %s", r.Code, r.Message)
}

// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc/codes"
)

// The refusal of a call that failed in the client, the same on HTTP and gRPC and in every
// Permguard SDK (BUG-0284): class ClassUnavailable, code CodeTimeout when the call's timeout or
// the context's deadline passed, CodeUnreachable when the data plane could not be reached.
// A caller's own cancellation is not a refusal: it is the context's error.
const (
	ClassUnavailable = "unavailable"
	CodeTimeout      = "timeout"
	CodeUnreachable  = "unreachable"
)

// clientFailure is the refusal of a call that failed in the client: a timeout, or a data plane
// that could not be reached. The caller's own cancellation is answered with the context's error.
// It answers nil when err is neither.
func clientFailure(ctx context.Context, err error) error {
	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(err, context.DeadlineExceeded) {
		return &Refusal{
			Class:   ClassUnavailable,
			Code:    CodeTimeout,
			Message: "the data plane did not answer within the call's timeout",
		}
	}
	var timeout interface{ Timeout() bool }
	if errors.As(err, &timeout) && timeout.Timeout() {
		return &Refusal{
			Class:   ClassUnavailable,
			Code:    CodeTimeout,
			Message: "the data plane did not answer within the call's timeout",
		}
	}
	var dial *net.OpError
	if errors.As(err, &dial) {
		return &Refusal{
			Class:   ClassUnavailable,
			Code:    CodeUnreachable,
			Message: fmt.Sprintf("the data plane could not be reached: %v", err),
		}
	}
	return nil
}

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

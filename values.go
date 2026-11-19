// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

// exactLimit is the largest integer both transports carry exactly, the server's bound: 2^53.
var exactLimit = new(big.Int).Lsh(big.NewInt(1), 53)

// checkValues refuses what the PDP refuses, before either transport sends it (BUG-0306,
// BUG-0307): an integer past ±2^53 is Refusal{validation, value_unrepresentable}; a value that is
// not JSON at all (NaN, an infinity, a channel, a function) is Refusal{validation,
// payload_malformed}. Both transports then send the same request.
func checkValues(request *EvaluateRequest) error {
	encoded, err := json.Marshal(request)
	if err != nil {
		return &Refusal{Class: "validation", Code: "payload_malformed", Message: fmt.Sprintf("the request is not JSON: %v", err)}
	}
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return &Refusal{Class: "validation", Code: "payload_malformed", Message: fmt.Sprintf("the request is not JSON: %v", err)}
	}
	return checkNumbers(decoded)
}

func checkNumbers(value any) error {
	switch held := value.(type) {
	case json.Number:
		text := held.String()
		if strings.ContainsAny(text, ".eE") {
			return nil
		}
		integer, ok := new(big.Int).SetString(text, 10)
		if !ok || new(big.Int).Abs(integer).Cmp(exactLimit) > 0 {
			return &Refusal{
				Class:   "validation",
				Code:    "value_unrepresentable",
				Message: fmt.Sprintf("%s is beyond the largest integer both transports carry exactly (2^53): send it as a string", text),
			}
		}
	case []any:
		for _, item := range held {
			if err := checkNumbers(item); err != nil {
				return err
			}
		}
	case map[string]any:
		for _, item := range held {
			if err := checkNumbers(item); err != nil {
				return err
			}
		}
	}
	return nil
}

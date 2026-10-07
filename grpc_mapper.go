// Copyright (c) 2022 Nitro Agility S.r.l.
// SPDX-License-Identifier: Apache-2.0

package permguard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	pdpv1 "github.com/permguard/sdk-go/internal/grpc/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

const maxExactProtoInteger = 1<<53 - 1

func requestToProto(request *EvaluateRequest) (*pdpv1.EvaluateRequest, error) {
	semantic, err := semanticToProto(request.Options)
	if err != nil {
		return nil, err
	}
	evaluations := make([]*pdpv1.Evaluation, 0, len(request.Evaluations))
	for index := range request.Evaluations {
		evaluation, err := evaluationToProto(&request.Evaluations[index])
		if err != nil {
			return nil, fmt.Errorf("evaluations[%d]: %w", index, err)
		}
		evaluations = append(evaluations, evaluation)
	}
	inputs, err := inputsToProto(request.PartitionInputs)
	if err != nil {
		return nil, fmt.Errorf("partition_inputs: %w", err)
	}
	context, err := mapToProto(request.Context)
	if err != nil {
		return nil, fmt.Errorf("context: %w", err)
	}
	subject, err := entityToProto(request.Subject)
	if err != nil {
		return nil, fmt.Errorf("subject: %w", err)
	}
	resource, err := entityToProto(request.Resource)
	if err != nil {
		return nil, fmt.Errorf("resource: %w", err)
	}
	action, err := actionToProto(request.Action)
	if err != nil {
		return nil, fmt.Errorf("action: %w", err)
	}
	principal, err := entityToProto(request.Principal)
	if err != nil {
		return nil, fmt.Errorf("principal: %w", err)
	}

	return &pdpv1.EvaluateRequest{
		Zone:                request.Zone,
		Ledger:              request.Ledger,
		Profile:             request.Profile,
		Subject:             subject,
		Resource:            resource,
		Action:              action,
		Context:             context,
		Principal:           principal,
		Evaluations:         evaluations,
		EvaluationsSemantic: semantic,
		RequestId:           request.RequestID,
		PartitionInputs:     inputs,
	}, nil
}

func semanticToProto(options *EvaluationOptions) (pdpv1.EvaluationsSemantic, error) {
	if options == nil || options.EvaluationsSemantic == "" {
		return pdpv1.EvaluationsSemantic_EVALUATIONS_SEMANTIC_UNSPECIFIED, nil
	}
	switch options.EvaluationsSemantic {
	case ExecuteAll:
		return pdpv1.EvaluationsSemantic_EVALUATIONS_SEMANTIC_EXECUTE_ALL, nil
	case DenyOnFirstDeny:
		return pdpv1.EvaluationsSemantic_EVALUATIONS_SEMANTIC_DENY_ON_FIRST_DENY, nil
	case PermitOnFirstPermit:
		return pdpv1.EvaluationsSemantic_EVALUATIONS_SEMANTIC_PERMIT_ON_FIRST_PERMIT, nil
	default:
		return 0, fmt.Errorf("unknown evaluations semantic %q", options.EvaluationsSemantic)
	}
}

func evaluationToProto(evaluation *Evaluation) (*pdpv1.Evaluation, error) {
	var context *structpb.Struct
	if evaluation.Context != nil {
		var err error
		context, err = mapToProto(*evaluation.Context)
		if err != nil {
			return nil, fmt.Errorf("context: %w", err)
		}
	}
	var inputs *pdpv1.PartitionInputs
	if evaluation.PartitionInputs != nil {
		mapped, err := inputsToProto(*evaluation.PartitionInputs)
		if err != nil {
			return nil, fmt.Errorf("partition_inputs: %w", err)
		}
		inputs = &pdpv1.PartitionInputs{Inputs: mapped}
	}
	subject, err := entityToProto(evaluation.Subject)
	if err != nil {
		return nil, fmt.Errorf("subject: %w", err)
	}
	resource, err := entityToProto(evaluation.Resource)
	if err != nil {
		return nil, fmt.Errorf("resource: %w", err)
	}
	action, err := actionToProto(evaluation.Action)
	if err != nil {
		return nil, fmt.Errorf("action: %w", err)
	}
	return &pdpv1.Evaluation{
		Subject:         subject,
		Resource:        resource,
		Action:          action,
		Context:         context,
		RequestId:       evaluation.RequestID,
		PartitionInputs: inputs,
	}, nil
}

func entityToProto(entity *Entity) (*pdpv1.Entity, error) {
	if entity == nil {
		return nil, nil
	}
	properties, err := mapToProto(entity.Properties)
	if err != nil {
		return nil, fmt.Errorf("properties: %w", err)
	}
	return &pdpv1.Entity{Type: entity.Type, Id: entity.ID, Properties: properties}, nil
}

func actionToProto(action *Action) (*pdpv1.Action, error) {
	if action == nil {
		return nil, nil
	}
	properties, err := mapToProto(action.Properties)
	if err != nil {
		return nil, fmt.Errorf("properties: %w", err)
	}
	return &pdpv1.Action{Name: action.Name, Properties: properties}, nil
}

func inputsToProto(inputs PartitionInputs) (map[string]*pdpv1.PartitionInput, error) {
	mapped := make(map[string]*pdpv1.PartitionInput, len(inputs))
	for name, input := range inputs {
		var data *structpb.Value
		if input.Data != nil {
			var err error
			data, err = valueToProto(input.Data)
			if err != nil {
				return nil, fmt.Errorf("%s.data: %w", name, err)
			}
		}
		mapped[name] = &pdpv1.PartitionInput{Type: input.Type, Data: data}
	}
	return mapped, nil
}

func mapToProto(values map[string]any) (*structpb.Struct, error) {
	if values == nil {
		return nil, nil
	}
	fields := make(map[string]*structpb.Value, len(values))
	for name, value := range values {
		mapped, err := valueToProto(value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		fields[name] = mapped
	}
	return &structpb.Struct{Fields: fields}, nil
}

func valueToProto(value any) (*structpb.Value, error) {
	if value == nil {
		return structpb.NewNullValue(), nil
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("value is not JSON: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return nil, fmt.Errorf("value is not JSON: %w", err)
	}
	normalized, err := normalizeJSON(decoded)
	if err != nil {
		return nil, err
	}
	return structpb.NewValue(normalized)
}

func normalizeJSON(value any) (any, error) {
	switch held := value.(type) {
	case json.Number:
		text := held.String()
		if !strings.ContainsAny(text, ".eE") {
			integer, err := strconv.ParseInt(text, 10, 64)
			if err != nil || integer > maxExactProtoInteger || integer < -maxExactProtoInteger {
				return nil, fmt.Errorf("integer %s is not exactly representable by protobuf Value", text)
			}
			return float64(integer), nil
		}
		number, err := strconv.ParseFloat(text, 64)
		if err != nil || math.IsInf(number, 0) || math.IsNaN(number) {
			return nil, fmt.Errorf("number %s is not representable by protobuf Value", text)
		}
		return number, nil
	case []any:
		mapped := make([]any, len(held))
		for index, item := range held {
			value, err := normalizeJSON(item)
			if err != nil {
				return nil, fmt.Errorf("item %d: %w", index, err)
			}
			mapped[index] = value
		}
		return mapped, nil
	case map[string]any:
		mapped := make(map[string]any, len(held))
		for name, item := range held {
			value, err := normalizeJSON(item)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
			mapped[name] = value
		}
		return mapped, nil
	default:
		return held, nil
	}
}

func responseFromProto(response *pdpv1.EvaluateResponse) *EvaluateResponse {
	if response == nil {
		return nil
	}
	evaluations := make([]Decision, 0, len(response.Evaluations))
	for _, evaluation := range response.Evaluations {
		evaluations = append(evaluations, decisionFromProto(evaluation))
	}
	return &EvaluateResponse{
		Decision:    response.Decision,
		RequestID:   response.RequestId,
		Context:     contextFromProto(response.Context),
		Evaluations: evaluations,
	}
}

func decisionFromProto(decision *pdpv1.Decision) Decision {
	if decision == nil {
		return Decision{}
	}
	return Decision{
		Decision:  decision.Decision,
		RequestID: decision.RequestId,
		Context:   contextFromProto(decision.Context),
	}
}

func contextFromProto(context *pdpv1.DecisionContext) *DecisionContext {
	if context == nil {
		return nil
	}
	return &DecisionContext{
		ID:           context.Id,
		ReasonAdmin:  reasonFromProto(context.ReasonAdmin),
		ReasonUser:   reasonFromProto(context.ReasonUser),
		Policies:     append([]string(nil), context.Policies...),
		AbsentInputs: append([]string(nil), context.AbsentInputs...),
	}
}

func reasonFromProto(reason *pdpv1.Reason) *Reason {
	if reason == nil {
		return nil
	}
	return &Reason{Code: reason.Code, Message: reason.Message}
}

func configurationFromProto(response *pdpv1.GetConfigurationResponse) *Configuration {
	if response == nil {
		return nil
	}
	configuration := &Configuration{
		Interface:    response.Interface,
		PDP:          response.Pdp,
		Capabilities: append([]string(nil), response.Capabilities...),
	}
	if response.Endpoints != nil {
		configuration.Endpoints = Endpoints{
			Evaluation:  response.Endpoints.Evaluation,
			Evaluations: response.Endpoints.Evaluations,
		}
	}
	if response.StoreScope != nil {
		configuration.StoreScope = StoreScope{
			In:      response.StoreScope.In,
			Zone:    response.StoreScope.Zone,
			Ledger:  response.StoreScope.Ledger,
			Profile: response.StoreScope.Profile,
		}
	}
	return configuration
}

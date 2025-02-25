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
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/permguard/permguard-go/az/azreq"
)

// Request

// MapPolicyStoreToGrpcPolicyStore maps the client policy store to the gRPC policy store.
func MapPolicyStoreToGrpcPolicyStore(policyStore *azreq.PolicyStore) (*PolicyStore, error) {
	if policyStore == nil {
		return nil, nil
	}
	target := &PolicyStore{}
	target.ID = policyStore.ID
	target.Kind = policyStore.Kind
	return target, nil
}

// MapPrincipalToGrpcPrincipal maps the client principal to the gRPC principal.
func MapPrincipalToGrpcPrincipal(principal *azreq.Principal) (*Principal, error) {
	if principal == nil {
		return nil, nil
	}
	target := &Principal{}
	target.ID = principal.ID
	target.Type = principal.Type
	source := principal.Source
	if source != "" {
		target.Source = &source
	}
	return target, nil
}

// MapEntitiesToGrpcEntities maps the client entities to the gRPC entities.
func MapEntitiesToGrpcEntities(entities *azreq.Entities) (*Entities, error) {
	if entities == nil {
		return nil, nil
	}
	target := &Entities{}
	target.Schema = entities.Schema
	originalItems := entities.Items
	if originalItems != nil {
		mappedItems := make([]*structpb.Struct, 0, len(originalItems))
		for _, item := range originalItems {
			data, err := structpb.NewStruct(item)
			if err != nil {
				return nil, err
			}
			mappedItems = append(mappedItems, data)
		}
		target.Items = mappedItems
	}
	return target, nil
}

// MapSubjectToGrpcSubject maps the client subject to the gRPC subject.
func MapSubjectToGrpcSubject(subject *azreq.Subject) (*Subject, error) {
	if subject == nil {
		return nil, nil
	}
	target := &Subject{}
	target.ID = subject.ID
	target.Type = subject.Type
	source := subject.Source
	if source != "" {
		target.Source = &source
	}
	properties := subject.Properties
	if properties != nil {
		data, err := structpb.NewStruct(properties)
		if err != nil {
			return nil, err
		}
		target.Properties = data
	}
	return target, nil
}

// MapResourceToGrpcResource maps the client resource to the gRPC resource.
func MapResourceToGrpcResource(resource *azreq.Resource) (*Resource, error) {
	if resource == nil {
		return nil, nil
	}
	target := &Resource{}
	target.ID = resource.ID
	target.Type = resource.Type
	properties := resource.Properties
	if properties != nil {
		data, err := structpb.NewStruct(properties)
		if err != nil {
			return nil, err
		}
		target.Properties = data
	}
	return target, nil
}

// MapActionToGrpcAction maps the client action to the gRPC action.
func MapActionToGrpcAction(action *azreq.Action) (*Action, error) {
	if action == nil {
		return nil, nil
	}
	target := &Action{}
	target.Name = action.Name
	properties := action.Properties
	if properties != nil {
		data, err := structpb.NewStruct(properties)
		if err != nil {
			return nil, err
		}
		target.Properties = data
	}
	return target, nil
}

// MapEvaluationToGrpcEvaluationRequest maps the client evaluation to the gRPC evaluation request.
func MapEvaluationToGrpcEvaluationRequest(evaluation *azreq.Evaluation) (*EvaluationRequest, error) {
	if evaluation == nil {
		return nil, nil
	}
	target := &EvaluationRequest{}
	requestID := evaluation.RequestID
	target.RequestID = &requestID
	subject := evaluation.Subject
	if subject != nil {
		subject, err := MapSubjectToGrpcSubject(subject)
		if err != nil {
			return nil, err
		}
		target.Subject = subject
	}
	resource := evaluation.Resource
	if resource != nil {
		resource, err := MapResourceToGrpcResource(resource)
		if err != nil {
			return nil, err
		}
		target.Resource = resource
	}
	action := evaluation.Action
	if action != nil {
		action, err := MapActionToGrpcAction(action)
		if err != nil {
			return nil, err
		}
		target.Action = action
	}
	ctx := evaluation.Context
	if ctx != nil {
		data, err := structpb.NewStruct(ctx)
		if err != nil {
			return nil, err
		}
		target.Context = data
	}
	return target, nil
}

// MapAuthZModelToGrpcAuthorizationModelRequest maps the client azrequest to the gRPC authorization model request.
func MapAuthZModelToGrpcAuthorizationModelRequest(azModel *azreq.AZModel) (*AuthorizationModelRequest, error) {
	req := &AuthorizationModelRequest{}
	req.ZoneID = int64(azModel.ZoneID)
	policyStore := azModel.PolicyStore
	if policyStore != nil {
		policyStore, err := MapPolicyStoreToGrpcPolicyStore(policyStore)
		if err != nil {
			return nil, err
		}
		req.PolicyStore = policyStore
	}
	principal := azModel.Principal
	if principal != nil {
		principal, err := MapPrincipalToGrpcPrincipal(principal)
		if err != nil {
			return nil, err
		}
		req.Principal = principal
	}
	entities := azModel.Entities
	if entities != nil {
		entities, err := MapEntitiesToGrpcEntities(entities)
		if err != nil {
			return nil, err
		}
		req.Entities = entities
	}
	return req, nil
}

// MapAZRequestToGrpcAuthorizationCheckRequest maps the client azrequest to the gRPC authorization check request.
func MapAZRequestToGrpcAuthorizationCheckRequest(azRequest *azreq.AZRequest) (*AuthorizationCheckRequest, error) {
	if azRequest == nil {
		return nil, nil
	}
	req := &AuthorizationCheckRequest{}
	azModel := azRequest.AZModel
	if azModel != nil {
		AuthorizationModel, err := MapAuthZModelToGrpcAuthorizationModelRequest(azModel)
		if err != nil {
			return nil, err
		}
		req.AuthorizationModel = AuthorizationModel
	}
	azEvaluations := azRequest.Evaluations
	if azEvaluations != nil {
		evaluations := []*EvaluationRequest{}
		for _, azEvaluation := range azEvaluations {
			evaluation, err := MapEvaluationToGrpcEvaluationRequest(&azEvaluation)
			if err != nil {
				return nil, err
			}
			if evaluation.RequestID == nil {
				requesID := evaluation.GetRequestID()
				evaluation.RequestID = &requesID
			}
			evaluations = append(evaluations, evaluation)
		}
		req.Evaluations = evaluations
	}
	return req, nil
}

// Response

// MapGrpcReasonResponseToReasonResponse maps the gRPC reason response to the reason response.
func MapGrpcReasonResponseToReasonResponse(reasonResponse *ReasonResponse) (*azreq.ReasonResponse, error) {
	if reasonResponse == nil {
		return nil, nil
	}
	target := &azreq.ReasonResponse{}
	target.Code = reasonResponse.Code
	target.Message = reasonResponse.Message
	return target, nil
}

// MapGrpcContextResponseToContextResponse maps the gRPC context response to the context response.
func MapGrpcContextResponseToContextResponse(contextResponse *ContextResponse) (*azreq.ContextResponse, error) {
	if contextResponse == nil {
		return nil, nil
	}
	target := &azreq.ContextResponse{}
	target.ID = contextResponse.ID
	if contextResponse.ReasonAdmin != nil {
		reasonAdmin, err := MapGrpcReasonResponseToReasonResponse(contextResponse.ReasonAdmin)
		if err != nil {
			return nil, err
		}
		target.ReasonAdmin = reasonAdmin
	}
	if contextResponse.ReasonUser != nil {
		reasonUser, err := MapGrpcReasonResponseToReasonResponse(contextResponse.ReasonUser)
		if err != nil {
			return nil, err
		}
		target.ReasonUser = reasonUser
	}
	return target, nil
}

// MapGrpcEvaluationResponseToEvaluationResponse maps the gRPC evaluation response to the evaluation response.
func MapGrpcEvaluationResponseToEvaluationResponse(evaluationResponse *EvaluationResponse) (*azreq.EvaluationResponse, error) {
	if evaluationResponse == nil {
		return nil, nil
	}
	target := &azreq.EvaluationResponse{}
	target.Decision = evaluationResponse.Decision
	if evaluationResponse.RequestID != nil {
		target.RequestID = *evaluationResponse.RequestID
	} else {
		target.RequestID = ""
	}
	if evaluationResponse.Context != nil {
		context, err := MapGrpcContextResponseToContextResponse(evaluationResponse.Context)
		if err != nil {
			return nil, err
		}
		target.Context = context
	}
	return target, nil
}

// MapGrpcAuthorizationCheckResponseToAZResponse maps the gRPC authorization check response to the authorization check response.
func MapGrpcAuthorizationCheckResponseToAZResponse(response *AuthorizationCheckResponse) (*azreq.AZResponse, error) {
	if response == nil {
		return nil, nil
	}
	target := &azreq.AZResponse{}
	target.Decision = response.Decision
	if response.RequestID != nil {
		target.RequestID = *response.RequestID
	} else {
		target.RequestID = ""
	}
	if response.Context != nil {
		context, err := MapGrpcContextResponseToContextResponse(response.Context)
		if err != nil {
			return nil, err
		}
		target.Context = context
	}
	if response.Evaluations != nil {
		evaluations := []azreq.EvaluationResponse{}
		for _, evaluationResponse := range response.Evaluations {
			evaluation, err := MapGrpcEvaluationResponseToEvaluationResponse(evaluationResponse)
			if err != nil {
				return nil, err
			}
			evaluations = append(evaluations, *evaluation)
		}
		target.Evaluations = evaluations
	}
	return target, nil
}

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

// MapPolicyStoreToGrpcPolicyStore maps the client policy store to the gRPC policy store.
func MapPolicyStoreToGrpcPolicyStore(policyStore *azreq.PolicyStore) (*PolicyStore, error) {
	if policyStore == nil {
		return nil, nil
	}
	target := &PolicyStore{}
	target.ID = policyStore.GetID()
	target.Kind = policyStore.GetKind()
	return target, nil
}

// MapPrincipalToGrpcPrincipal maps the client principal to the gRPC principal.
func MapPrincipalToGrpcPrincipal(principal *azreq.Principal) (*Principal, error) {
	if principal == nil {
		return nil, nil
	}
	target := &Principal{}
	target.ID = principal.GetID()
	target.Kind = principal.GetKind()
	source := principal.GetSource()
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
	target.Schema = entities.GetSchema()
	originalItems := entities.GetItems()
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
	target.ID = subject.GetID()
	target.Kind = subject.GetKind()
	source := subject.GetSource()
	if source != "" {
		target.Source = &source
	}
	properties := subject.GetProperties()
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
	target.ID = resource.GetID()
	target.Kind = resource.GetKind()
	properties := resource.GetProperties()
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
	target.Name = action.GetName()
	properties := action.GetProperties()
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
	requestID := evaluation.GetRequestID()
	target.RequestID = &requestID
	subject := evaluation.GetSubject()
	if subject != nil {
		subject, err := MapSubjectToGrpcSubject(subject)
		if err != nil {
			return nil, err
		}
		target.Subject = subject
	}
	resource := evaluation.GetResource()
	if resource != nil {
		resource, err := MapResourceToGrpcResource(resource)
		if err != nil {
			return nil, err
		}
		target.Resource = resource
	}
	action := evaluation.GetAction()
	if action != nil {
		action, err := MapActionToGrpcAction(action)
		if err != nil {
			return nil, err
		}
		target.Action = action
	}
	ctx := evaluation.GetContext()
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
func MapAuthZModelToGrpcAuthorizationModelRequest(azModel *azreq.AuthZModel) (*AuthorizationModelRequest, error) {
	req := &AuthorizationModelRequest{}
	req.ZoneID = int64(azModel.GetZoneID())
	policyStore := azModel.GetPolicyStore()
	if policyStore != nil {
		policyStore, err := MapPolicyStoreToGrpcPolicyStore(policyStore)
		if err != nil {
			return nil, err
		}
		req.PolicyStore = policyStore
	}
	principal := azModel.GetPrincipal()
	if principal != nil {
		principal, err := MapPrincipalToGrpcPrincipal(principal)
		if err != nil {
			return nil, err
		}
		req.Principal = principal
	}
	entities := azModel.GetEntities()
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
	azModel := azRequest.GetAuthZModel()
	if azModel != nil {
		AuthorizationModel, err := MapAuthZModelToGrpcAuthorizationModelRequest(azModel)
		if err != nil {
			return nil, err
		}
		req.AuthorizationModel = AuthorizationModel
	}
	azEvaluations := azRequest.GetEvaluations()
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

// // MapAgentReasonResponseToGrpcReasonResponse maps the agent reason response to the gRPC reason response.
// func MapAgentReasonResponseToGrpcReasonResponse(reasonResponse *azreq.ReasonResponse) (*ReasonResponse, error) {
// 	if reasonResponse == nil {
// 		return nil, nil
// 	}
// 	target := &ReasonResponse{}
// 	target.Code = reasonResponse.Code
// 	target.Message = reasonResponse.Message
// 	return target, nil
// }

// // MapAgentContextResponseToGrpcContextResponse maps the agent context response to the gRPC context response.
// func MapAgentContextResponseToGrpcContextResponse(contextResponse *azreq.ContextResponse) (*ContextResponse, error) {
// 	if contextResponse == nil {
// 		return nil, nil
// 	}
// 	target := &ContextResponse{}
// 	target.ID = contextResponse.ID
// 	if contextResponse.ReasonAdmin != nil {
// 		reasonAdmin, err := MapAgentReasonResponseToGrpcReasonResponse(contextResponse.ReasonAdmin)
// 		if err != nil {
// 			return nil, err
// 		}
// 		target.ReasonAdmin = reasonAdmin
// 	}
// 	if contextResponse.ReasonUser != nil {
// 		reasonUser, err := MapAgentReasonResponseToGrpcReasonResponse(contextResponse.ReasonUser)
// 		if err != nil {
// 			return nil, err
// 		}
// 		target.ReasonUser = reasonUser
// 	}
// 	return target, nil
// }

// // MapAgentEvaluationResponseToGrpcEvaluationResponse maps the agent evaluation response to the gRPC evaluation response.
// func MapAgentEvaluationResponseToGrpcEvaluationResponse(evaluationResponse *azreq.EvaluationResponse) (*EvaluationResponse, error) {
// 	if evaluationResponse == nil {
// 		return nil, nil
// 	}
// 	target := &EvaluationResponse{}
// 	target.Decision = evaluationResponse.Decision
// 	target.RequestID = &evaluationResponse.RequestID
// 	if evaluationResponse.Context != nil {
// 		context, err := MapAgentContextResponseToGrpcContextResponse(evaluationResponse.Context)
// 		if err != nil {
// 			return nil, err
// 		}
// 		target.Context = context
// 	}
// 	return target, nil
// }

// // MapAgentAuthorizationCheckResponseToGrpcAuthorizationCheckResponse maps the agent authorization check response to the gRPC authorization check response.
// func MapAgentAuthorizationCheckResponseToGrpcAuthorizationCheckResponse(response *azreq.AuthorizationCheckResponse) (*AuthorizationCheckResponse, error) {
// 	if response == nil {
// 		return nil, nil
// 	}
// 	target := &AuthorizationCheckResponse{}
// 	target.RequestID = &response.RequestID
// 	target.Decision = response.Decision
// 	if response.Context != nil {
// 		context, err := MapAgentContextResponseToGrpcContextResponse(response.Context)
// 		if err != nil {
// 			return nil, err
// 		}
// 		target.Context = context
// 	}
// 	if response.Evaluations != nil {
// 		evaluations := []*EvaluationResponse{}
// 		for _, evaluationResponse := range response.Evaluations {
// 			evaluation, err := MapAgentEvaluationResponseToGrpcEvaluationResponse(&evaluationResponse)
// 			if err != nil {
// 				return nil, err
// 			}
// 			evaluations = append(evaluations, evaluation)
// 		}
// 		target.Evaluations = evaluations
// 	}
// 	return target, nil
// }

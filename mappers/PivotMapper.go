package mappers

import (
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/models"
)

func ToAgentOperatorModel(req *pivot.AssignToOpRequest) *models.AgentOperator {
	if req == nil || req.AgentId == nil {
		return nil
	}
	return &models.AgentOperator{
		OperatorID: req.OperatorId,
		AgentID:    *req.AgentId,
	}
}

func ToArrAgentOperatorModel(req *pivot.ArraysToOpRequest) *[]models.AgentOperator {
	if req == nil || req.AgentIds == nil {
		return nil
	}

	arrayModel := make([]models.AgentOperator, 0, len(*req.AgentIds))

	for _, agentId := range *req.AgentIds {
		arrayModel = append(arrayModel, models.AgentOperator{
			OperatorID: req.OperatorId,
			AgentID:    agentId,
		})
	}

	return &arrayModel
}

func ToAgentOperatorResponse(model *models.AgentOperator) *pivot.AssignToOpResponse {
	if model == nil {
		return nil
	}

	agentID := model.AgentID

	return &pivot.AssignToOpResponse{
		OperatorID: model.OperatorID,
		AgentID:    &agentID,
	}
}

func ToArrAgentOperatorResponse(agentOperators *[]models.AgentOperator) *pivot.ArraysToOpResponse {
	if agentOperators == nil || len(*agentOperators) == 0 {
		return nil
	}

	agentIds := make([]uint, 0, len(*agentOperators))
	for _, ao := range *agentOperators {
		agentIds = append(agentIds, ao.AgentID)
	}

	return &pivot.ArraysToOpResponse{
		OperatorId: (*agentOperators)[0].OperatorID,
		AgentIds:   &agentIds,
	}
}

func ToAgencyOperatorModel(req *pivot.AssignToOpRequest) *models.AgencyOperator {
	if req == nil || req.AgencyId == nil {
		return nil
	}
	return &models.AgencyOperator{
		OperatorID: req.OperatorId,
		AgencyID:   *req.AgencyId,
	}
}

func ToArrAgencyOperatorModel(req *pivot.ArraysToOpRequest) *[]models.AgencyOperator {
	if req == nil || req.AgencyIds == nil {
		return nil
	}

	arrayModel := make([]models.AgencyOperator, 0, len(*req.AgencyIds))

	for _, agencyId := range *req.AgencyIds {
		arrayModel = append(arrayModel, models.AgencyOperator{
			OperatorID: req.OperatorId,
			AgencyID:   agencyId,
		})
	}

	return &arrayModel
}

func ToAgencyOperatorResponse(model *models.AgencyOperator) *pivot.AssignToOpResponse {
	if model == nil {
		return nil
	}

	agencyID := model.AgencyID

	return &pivot.AssignToOpResponse{
		OperatorID: model.OperatorID,
		AgencyID:   &agencyID,
	}
}

func ToArrAgencyOperatorResponse(agencyOperators *[]models.AgencyOperator) *pivot.ArraysToOpResponse {
	if agencyOperators == nil || len(*agencyOperators) == 0 {
		return nil
	}

	agencyIds := make([]uint, 0, len(*agencyOperators))
	for _, ao := range *agencyOperators {
		agencyIds = append(agencyIds, ao.AgencyID)
	}

	return &pivot.ArraysToOpResponse{
		OperatorId: (*agencyOperators)[0].OperatorID,
		AgencyIds:  &agencyIds,
	}
}

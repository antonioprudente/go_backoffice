package mappers

import (
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/models"
)

func ToAgentOperatorModel(req *pivot.AssignToOpRequest) *models.AgentOperator {
	model := &models.AgentOperator{
		OperatorID: req.OperatorId,
		AgentID:    *req.AgentId,
	}

	return model
}

func ToArrAgentOperatorModel(req *pivot.ArraysToOpRequest) *[]models.AgentOperator {
	if req == nil {
		return nil
	}

	arrayModel := make([]models.AgentOperator, len(*req.AgentIds))

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

	res := &pivot.AssignToOpResponse{
		Operator: ToUserResponse(model.Operator),
		Agent:    ToUserResponse(model.Agent),
	}

	return res
}

func ToAgencyOperatorModel(req *pivot.AssignToOpRequest) *models.AgencyOperator {
	model := &models.AgencyOperator{
		OperatorID: req.OperatorId,
		AgencyID:   *req.AgencyId,
	}

	return model
}

func ToArrAgencyOperatorModel(req *pivot.ArraysToOpRequest) *[]models.AgencyOperator {
	if req == nil {
		return nil
	}

	arrayModel := make([]models.AgencyOperator, len(*req.AgencyIds))

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

	res := &pivot.AssignToOpResponse{
		Operator: ToUserResponse(model.Operator),
		Agency:   ToUserResponse(model.Agency),
	}

	return res
}

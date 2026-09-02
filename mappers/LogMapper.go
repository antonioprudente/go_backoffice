package mappers

import (
	"example/go_backoffice/dto/log"
	"example/go_backoffice/models"
)

func ToLogResponse(model *models.ActivityLog) *log.LogResponse {
	if model == nil {
		return nil
	}

	response := &log.LogResponse{
		ID:          model.ID,
		Action:      model.Action.String(),
		TargetType:  model.TargetType,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
	}

	var checkIds = model.ActorID != nil && model.TargetID != nil
	var checkEnt = model.Actor != nil

	if checkIds && !checkEnt {
		response.ActorID = model.ActorID
		response.TargetID = model.TargetID
	}

	if checkEnt {
		actorResp := ToUserResponse(model.Actor)
		response.Actor = &actorResp
		response.TargetID = model.TargetID
	}

	return response
}

func ToLogResponses(logs []*models.ActivityLog) []*log.LogResponse {
	res := make([]*log.LogResponse, len(logs))
	for i := range logs {
		res[i] = ToLogResponse(logs[i])
	}
	return res
}

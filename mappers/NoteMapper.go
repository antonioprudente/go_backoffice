package mappers

import (
	"example/go_backoffice/dto/note"
	"example/go_backoffice/models"
)

func ToNoteModel(req *note.NoteRequest) *models.Note {
	if req == nil {
		return nil
	}
	return &models.Note{
		ActorID:  &req.ActorID,
		TargetID: &req.TargetID,
		Content:  req.Content,
	}
}

func ToNoteResponse(model *models.Note, isFull bool) *note.NoteResponse {
	if model == nil {
		return nil
	}

	response := &note.NoteResponse{
		ID:        model.ID,
		Content:   model.Content,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	if !isFull {
		if model.Actor != nil {
			response.ActorUsername = &model.Actor.Username
		}
		if model.Target != nil {
			response.TargetUsername = &model.Target.Username
		}
		return response
	}

	if model.Actor != nil {
		actorResp := ToUserResponse(model.Actor)
		response.Actor = &actorResp
	}
	if model.Target != nil {
		targetResp := ToUserResponse(model.Target)
		response.Target = &targetResp
	}

	return response
}

func ToNoteResponses(notes []*models.Note) []*note.NoteResponse {
	res := make([]*note.NoteResponse, len(notes))
	for i := range notes {
		res[i] = ToNoteResponse(notes[i], false)
	}
	return res
}

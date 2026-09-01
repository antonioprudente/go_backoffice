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
		ActorID:  req.ActorID,
		TargetID: req.TargetID,
		Content:  req.Content,
	}
}

func ToNoteResponse(model *models.Note) *note.NoteResponse {
	if model == nil {
		return nil
	}
	return &note.NoteResponse{
		ID:        model.ID,
		ActorID:   model.ActorID,
		TargetID:  model.TargetID,
		Content:   model.Content,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func ToNoteResponses(notes []*models.Note) []*note.NoteResponse {
	res := make([]*note.NoteResponse, len(notes))
	for i := range notes {
		res[i] = ToNoteResponse(notes[i])
	}
	return res
}

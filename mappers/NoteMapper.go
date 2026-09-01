// mappers/NoteMapper.go
package mappers

import (
	"example/go_backoffice/dto/note"
	"example/go_backoffice/models"
)

func ToNoteModel(req *note.NoteRequest) *models.Note {
	return &models.Note{
		ActorID:  req.ActorID,
		TargetID: req.TargetID,
		Content:  req.Content,
	}
}

func ToNoteResponse(model *models.Note) *note.NoteResponse {
	return &note.NoteResponse{
		ID:        model.ID,
		ActorID:   model.ActorID,
		TargetID:  model.TargetID,
		Content:   model.Content,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

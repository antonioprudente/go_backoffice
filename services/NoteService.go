package services

import (
	"example/go_backoffice/dto/note"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/models"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"
	"fmt"
	"reflect"
)

type NoteService interface {
	AssignNote(req *note.NoteRequest, actor policies.AuthContext) (*note.NoteResponse, error)
}

type noteService struct {
	repo       repositories.NoteRepo
	userRepo   repositories.UserRepo
	logService ActivityLogService
	policy     policies.UserPolicy
}

func NewNoteService(
	repo repositories.NoteRepo,
	userRepo repositories.UserRepo,
	logService ActivityLogService,
	policy policies.UserPolicy,
) NoteService {
	return &noteService{
		repo:       repo,
		userRepo:   userRepo,
		logService: logService,
		policy:     policy,
	}
}

func (s *noteService) AssignNote(req *note.NoteRequest, actor policies.AuthContext) (*note.NoteResponse, error) {
	model := mappers.ToNoteModel(req)

	targetUser, err := s.userRepo.GetByID(req.TargetID)
	if err != nil {
		return nil, err
	}

	if err := s.policy.View(actor, targetUser); err != nil {
		return nil, err
	}

	newNote, err := s.repo.NewNote(*model)
	if err != nil {
		return nil, err
	}

	err = s.logService.NewLog(&models.ActivityLog{
		ActorID:     &actor.UserID,
		ActorRole:   enums.Role(actor.Role),
		Action:      enums.Create,
		TargetType:  reflect.TypeOf(targetUser).Elem().Name(),
		TargetID:    &targetUser.ID,
		Description: fmt.Sprintf("Creata una nota sull'utente '%s' con ruolo %s", targetUser.Username, targetUser.Role),
	})
	if err != nil {
		return nil, err
	}

	res := mappers.ToNoteResponse(newNote)
	return res, nil
}

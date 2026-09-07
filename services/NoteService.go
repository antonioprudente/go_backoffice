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

	"gorm.io/gorm"
)

type NoteService interface {
	AssignNote(req *note.NoteRequest, actor policies.AuthContext) (*note.NoteResponse, error)
	GetNoteByID(id uint, actor policies.AuthContext) (*note.NoteResponse, error)
	GetAllNotes(actor policies.AuthContext) ([]*note.NoteResponse, error)
	UpdateNote(id uint, req *note.NoteRequest, actor policies.AuthContext) (*note.NoteResponse, error)
	DeleteNote(id uint, actor policies.AuthContext) error
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
		ActorID: &actor.UserID,
		//ActorRole:   enums.Role(actor.Role),
		Action:      enums.Create,
		TargetType:  reflect.TypeOf(newNote).Elem().Name(),
		TargetID:    &newNote.ID,
		Description: fmt.Sprintf("L'utente %d ha creato una nota all'utente %d", actor.UserID, targetUser.ID),
	})
	if err != nil {
		return nil, err
	}

	res := mappers.ToNoteResponse(newNote)
	return res, nil
}

func (s *noteService) GetNoteByID(id uint, actor policies.AuthContext) (*note.NoteResponse, error) {
	note, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	target, err := s.userRepo.GetByID(*note.TargetID)
	if err != nil {
		return nil, err
	}

	if err := s.policy.View(actor, target); err != nil {
		return nil, err
	}

	response := mappers.ToNoteResponse(note)
	return response, nil
}

func (s *noteService) GetAllNotes(actor policies.AuthContext) ([]*note.NoteResponse, error) {
	var notes []*models.Note
	var err error

	if actor.Role == enums.RoleAdmin.String() {
		notes, err = s.repo.GetAll()
	} else {
		notes, err = s.repo.GetAllByActorID(actor.UserID)
	}

	if err != nil {
		return nil, err
	}

	return mappers.ToNoteResponses(notes), nil
}

func (s *noteService) UpdateNote(id uint, req *note.NoteRequest, actor policies.AuthContext) (*note.NoteResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	target, err := s.userRepo.GetByID(*existing.TargetID)
	if err != nil {
		return nil, err
	}

	if err := s.policy.View(actor, target); err != nil {
		return nil, err
	}

	if req.Content != "" || existing.Content != req.Content {
		existing.Content = req.Content
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	err = s.logService.NewLog(&models.ActivityLog{
		ActorID: &actor.UserID,
		//ActorRole:   enums.Role(actor.Role),
		Action:      enums.Update,
		TargetType:  reflect.TypeOf(existing).Elem().Name(),
		TargetID:    &existing.ID,
		Description: fmt.Sprintf("Aggiornata la nota %d associata all'utente %d", existing.ID, *existing.TargetID),
	})
	if err != nil {
		return nil, err
	}

	return mappers.ToNoteResponse(existing), nil
}

func (s *noteService) DeleteNote(id uint, actor policies.AuthContext) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	target, err := s.userRepo.GetByID(*existing.TargetID)
	if err != nil {
		return err
	}

	if err := s.policy.View(actor, target); err != nil {
		return err
	}

	deleted, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if !deleted {
		return gorm.ErrRecordNotFound
	}

	err = s.logService.NewLog(&models.ActivityLog{
		ActorID: &actor.UserID,
		//ActorRole:   enums.Role(actor.Role),
		Action:      enums.Delete,
		TargetType:  reflect.TypeOf(existing).Elem().Name(),
		TargetID:    &id,
		Description: fmt.Sprintf("Eliminata la nota %d associata all'utente %d", existing.ID, existing.TargetID),
	})
	if err != nil {
		return err
	}

	return nil
}

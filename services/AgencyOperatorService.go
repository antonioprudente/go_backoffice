package services

import (
	"errors"
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/models"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"
	"fmt"
)

type AgencyOperatorService interface {
	AssignAgencyToOperator(request *pivot.AssignToOpRequest, actor policies.AuthContext) (*pivot.AssignToOpResponse, error)
	AssignAgenciesToOperator(request *pivot.ArraysToOpRequest, actor policies.AuthContext) (*pivot.ArraysToOpResponse, error)
	RemoveAgencyFromOperator(agencyID *uint, operatorID *uint, actor policies.AuthContext) (*bool, error)
}

type agencyOperatorService struct {
	repo       repositories.AgencyOperatorRepo
	userRepo   repositories.UserRepo
	logService ActivityLogService
}

func NewAgencyOperatorService(
	repo repositories.AgencyOperatorRepo,
	userRepo repositories.UserRepo,
	logService ActivityLogService,
) AgencyOperatorService {
	return &agencyOperatorService{
		repo:       repo,
		userRepo:   userRepo,
		logService: logService,
	}
}

func (s *agencyOperatorService) AssignAgencyToOperator(request *pivot.AssignToOpRequest, actor policies.AuthContext) (*pivot.AssignToOpResponse, error) {
	if request == nil || request.AgencyId == nil {
		return nil, errors.New("agency_id mancante nella request")
	}

	_, err := s.userRepo.GetByIDAndRole(*request.AgencyId, enums.RoleAgency.String())
	if err != nil {
		return nil, errors.New("L'utente selezionato non è un'agenzia")
	}

	_, err = s.userRepo.GetByIDAndRole(request.OperatorId, enums.RoleOperator.String())
	if err != nil {
		return nil, errors.New("L'utente selezionato non è un operatore")
	}
	model := mappers.ToAgencyOperatorModel(request)

	newPivot, err := s.repo.AssignAgency(model)
	if err != nil {
		return nil, err
	}

	// Log dell'assegnazione singola
	_ = s.logService.NewLog(&models.ActivityLog{
		ActorID:     &actor.UserID,
		ActorRole:   enums.Role(actor.Role),
		Action:      enums.Assignment,
		TargetType:  "AgencyOperator",
		TargetID:    request.AgencyId,
		Description: fmt.Sprintf("Assegnata Agenzia #%d all'Operatore #%d", *request.AgencyId, request.OperatorId),
	})

	response := mappers.ToAgencyOperatorResponse(newPivot)
	return response, nil
}

func (s *agencyOperatorService) AssignAgenciesToOperator(request *pivot.ArraysToOpRequest, actor policies.AuthContext) (*pivot.ArraysToOpResponse, error) {
	if request == nil {
		return nil, errors.New("request payload non valido")
	}
	model := mappers.ToArrAgencyOperatorModel(request)

	newPivots, err := s.repo.AssignAgenciesMassive(model)
	if err != nil {
		return nil, err
	}

	// Log dell'assegnazione massiva
	_ = s.logService.NewLog(&models.ActivityLog{
		ActorID:     &actor.UserID,
		ActorRole:   enums.Role(actor.Role),
		Action:      enums.Assignment,
		TargetType:  "AgencyOperator",
		TargetID:    &request.OperatorId,
		Description: fmt.Sprintf("Assegnate %d Agenzie all'Operatore #%d", len(*request.AgencyIds), request.OperatorId),
	})

	response := mappers.ToArrAgencyOperatorResponse(newPivots)
	return response, nil
}

func (s *agencyOperatorService) RemoveAgencyFromOperator(agencyID *uint, operatorID *uint, actor policies.AuthContext) (*bool, error) {
	res, err := s.repo.DeleteByAgencyIDAndOperatorID(*agencyID, *operatorID)
	if err != nil {
		return nil, err
	}

	if !res {
		return nil, errors.New("associazione tra agente e operatore non trovata")
	}

	// Log della rimozione
	_ = s.logService.NewLog(&models.ActivityLog{
		ActorID:     &actor.UserID,
		ActorRole:   enums.Role(actor.Role),
		Action:      enums.Remove,
		TargetType:  "AgencyOperator",
		TargetID:    agencyID,
		Description: fmt.Sprintf("Rimossa Agenzia #%d dall'Operatore #%d", *agencyID, *operatorID),
	})

	return &res, nil
}

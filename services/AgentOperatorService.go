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

type AgentOperatorService interface {
	AssignAgentToOperator(request *pivot.AssignToOpRequest, actor policies.AuthContext) (*pivot.AssignToOpResponse, error)
	AssignAgentsToOperator(request *pivot.ArraysToOpRequest, actor policies.AuthContext) (*pivot.ArraysToOpResponse, error)
	RemoveAgentFromOperator(agentID *uint, operatorID *uint, actor policies.AuthContext) (*bool, error)
}

type agentOperatorService struct {
	repo       repositories.AgentOperatorRepo
	userRepo   repositories.UserRepo
	logService ActivityLogService
}

func NewAgentOperatorService(
	repo repositories.AgentOperatorRepo,
	userRepo repositories.UserRepo,
	logService ActivityLogService,
) AgentOperatorService {
	return &agentOperatorService{
		repo:       repo,
		userRepo:   userRepo,
		logService: logService,
	}
}

func (s *agentOperatorService) AssignAgentToOperator(request *pivot.AssignToOpRequest, actor policies.AuthContext) (*pivot.AssignToOpResponse, error) {
	if actor.Role != enums.RoleAdmin.String() {
		return nil, errors.New("unauthorized: solo gli amministratori possono eseguire questa operazione")
	}

	if request == nil {
		return nil, errors.New("request payload non valido")
	}

	_, err := s.userRepo.GetByIDAndRole(*request.AgentId, enums.RoleAgent.String())
	if err != nil {
		return nil, errors.New("L'utente selezionato non è un agente")
	}

	_, err = s.userRepo.GetByIDAndRole(request.OperatorId, enums.RoleOperator.String())
	if err != nil {
		return nil, errors.New("L'utente selezionato non è un operatore")
	}

	model := mappers.ToAgentOperatorModel(request)

	newPivot, err := s.repo.AssignAgent(model)
	if err != nil {
		return nil, err
	}

	// Log dell'assegnazione
	_ = s.logService.NewLog(&models.ActivityLog{
		ActorID:     &actor.UserID,
		ActorRole:   enums.Role(actor.Role),
		Action:      "ASSIGN_AGENT_OPERATOR",
		TargetType:  "AgentOperator",
		TargetID:    request.AgentId,
		Description: fmt.Sprintf("Assegnato Agente #%d all'Operatore #%d", *request.AgentId, request.OperatorId),
	})

	response := mappers.ToAgentOperatorResponse(newPivot)
	return response, nil
}

func (s *agentOperatorService) AssignAgentsToOperator(request *pivot.ArraysToOpRequest, actor policies.AuthContext) (*pivot.ArraysToOpResponse, error) {
	if request == nil {
		return nil, errors.New("request payload non valido")
	}
	model := mappers.ToArrAgentOperatorModel(request)

	newPivots, err := s.repo.AssignAgentsMassive(model)
	if err != nil {
		return nil, err
	}

	// Log dell'assegnazione massiva
	_ = s.logService.NewLog(&models.ActivityLog{
		ActorID:     &actor.UserID,
		ActorRole:   enums.Role(actor.Role),
		Action:      "MASSIVE_ASSIGN_AGENT_OPERATOR",
		TargetType:  "AgentOperator",
		TargetID:    &request.OperatorId,
		Description: fmt.Sprintf("Assegnati %d Agenti all'Operatore #%d", len(*request.AgentIds), request.OperatorId),
	})

	response := mappers.ToArrAgentOperatorResponse(newPivots)
	return response, nil
}

func (s *agentOperatorService) RemoveAgentFromOperator(agentID *uint, operatorID *uint, actor policies.AuthContext) (*bool, error) {
	res, err := s.repo.DeleteByAgentIDAndOperatorID(*agentID, *operatorID)
	if err != nil {
		return nil, err
	}

	if !res {
		return nil, errors.New("associazione tra agente e operatore non trovata")
	}

	// Log della disassociazione
	_ = s.logService.NewLog(&models.ActivityLog{
		ActorID:     &actor.UserID,
		ActorRole:   enums.Role(actor.Role),
		Action:      "REMOVE_AGENT_OPERATOR",
		TargetType:  "AgentOperator",
		TargetID:    agentID,
		Description: fmt.Sprintf("Rimosso Agente #%d dall'Operatore #%d", *agentID, *operatorID),
	})

	return &res, nil
}

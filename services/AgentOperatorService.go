package services

import (
	"errors"
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"
)

type AgentOperatorService interface {
	AssignAgentToOperator(request *pivot.AssignToOpRequest, actor policies.AuthContext) (*pivot.AssignToOpResponse, error)
	AssignAgentsToOperator(request *pivot.ArraysToOpRequest, actor policies.AuthContext) (*pivot.ArraysToOpResponse, error)
	RemoveAgentFromOperator(agentID *uint, operatorID *uint, actor policies.AuthContext) (*bool, error)
}

type agentOperatorService struct {
	repo     repositories.AgentOperatorRepo
	userRepo repositories.UserRepo
}

func NewAgentOperatorService(
	repo repositories.AgentOperatorRepo,
	userRepo repositories.UserRepo,
) AgentOperatorService {
	return &agentOperatorService{
		repo:     repo,
		userRepo: userRepo,
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

	response := mappers.ToArrAgentOperatorResponse(newPivots)
	return response, nil
}

func (s *agentOperatorService) RemoveAgentFromOperator(agentID *uint, operatorID *uint, actor policies.AuthContext) (*bool, error) {
	res, err := s.repo.DeleteByAgentIDAndOperatorID(*agentID, *operatorID)
	if err != nil {
		return nil, err
	}

	// Se nessuna riga è stata cancellata, restituiamo un errore esplicito
	if !res {
		return nil, errors.New("associazione tra agente e operatore non trovata")
	}

	return &res, nil
}

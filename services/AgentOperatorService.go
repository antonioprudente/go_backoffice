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
}

type agentOperatorService struct {
	repo repositories.AgentOperatorRepo
}

func NewAgentOperatorService(repo repositories.AgentOperatorRepo) AgentOperatorService {
	return &agentOperatorService{repo: repo}
}

func (s *agentOperatorService) Create(request *pivot.AssignToOpRequest, actor policies.AuthContext) (*pivot.AssignToOpResponse, error) {
	if actor.Role != enums.RoleAdmin.String() {
		return nil, errors.New("unauthorized: solo gli amministratori possono eseguire questa operazione")
	}

	if request == nil {
		return nil, errors.New("request payload non valido")
	}

	model := mappers.ToAgentOperatorModel(request)

	newPivot, err := s.repo.AssignAgent(model)
	if err != nil {
		return nil, err
	}

	response := mappers.ToAgentOperatorResponse(newPivot)

	return response, nil

}

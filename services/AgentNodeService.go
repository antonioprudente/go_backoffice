// services/AgentNodeService.go
package services

import (
	"example/go_backoffice/dto/agent_node"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AgentNodeService interface {
	CreateNode(request *agent_node.AgentNodeRequest) (*agent_node.AgentNodeResponse, error)
}

type agentNodeService struct {
	repo        repositories.AgentNodeRepo
	userService UserService
}

func NewAgentNodeService(repo repositories.AgentNodeRepo, userService UserService) AgentNodeService {
	return &agentNodeService{
		repo:        repo,
		userService: userService,
	}
}

func (s *agentNodeService) CreateNode(request *agent_node.AgentNodeRequest) (*agent_node.AgentNodeResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Agent.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	request.Agent.Password = string(hashed)
	newNode := mappers.ToAgentNodeModel(request)
	newNode.Agent.Role = enums.RoleAgent

	if err := s.repo.Create(newNode); err != nil {
		return nil, err
	}

	response := mappers.ToAgentNodeResponse(newNode)
	return &response, nil
}

// services/AgentNodeService.go
package services

import (
	"example/go_backoffice/dto/agent_node"
	"example/go_backoffice/mappers"
	"example/go_backoffice/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AgentNodeService interface {
	CreateNode(request *agent_node.AgentNodeRequest) (*agent_node.AgentNodeResponse, error)
	GetTrees() ([]*agent_node.AgentNodeResponse, error)
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
	if request.Agent != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(request.Agent.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		request.Agent.Password = string(hashed)
	}

	newNode := mappers.ToAgentNodeModel(request)

	if err := s.repo.Create(newNode); err != nil {
		return nil, err
	}

	response := mappers.ToAgentNodeResponse(newNode)
	return &response, nil
}

func (s *agentNodeService) GetTrees() ([]*agent_node.AgentNodeResponse, error) {
	// 1. Recupera l'albero di modelli (slice di puntatori)
	treeModels, err := s.repo.GetTrees()
	if err != nil {
		return nil, err
	}

	// 2. Mappa l'intero albero in un'unica riga pulita
	return mappers.ToAgentNodePtrResponses(treeModels), nil
}

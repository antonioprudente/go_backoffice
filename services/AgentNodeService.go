package services

import (
	"example/go_backoffice/dto/agent_node"
	"example/go_backoffice/mappers"
	"example/go_backoffice/models"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AgentNodeService interface {
	CreateNode(request *agent_node.AgentNodeRequest, actor policies.AuthContext) (*agent_node.AgentNodeResponse, error)
	GetTrees() ([]*agent_node.AgentNodeResponse, error)
}

type agentNodeService struct {
	repo      repositories.AgentNodeRepo
	scopeRepo repositories.ScopeRepo
	policy    policies.UserPolicy
}

func NewAgentNodeService(
	repo repositories.AgentNodeRepo,
	scopeRepo repositories.ScopeRepo,
	policy policies.UserPolicy,
) AgentNodeService {
	return &agentNodeService{
		repo:      repo,
		scopeRepo: scopeRepo,
		policy:    policy,
	}
}

func (s *agentNodeService) CreateNode(request *agent_node.AgentNodeRequest, actor policies.AuthContext) (*agent_node.AgentNodeResponse, error) {
	if request.Agent != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(request.Agent.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		request.Agent.Password = string(hashed)
	}

	newNode := mappers.ToAgentNodeModel(request)
	if err := s.policy.Create(actor, newNode.Agent); err != nil {
		return nil, err
	}

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

func toSet(ids []uint) map[uint]bool {
	set := make(map[uint]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	return set
}

func filterRootsByAgentID(nodes []*models.AgentNode, allowed map[uint]bool) []*models.AgentNode {
	var result []*models.AgentNode
	for _, n := range nodes {
		filterChildren(n, allowed)
		if allowed[n.AgentID] {
			result = append(result, n)
		} else {
			result = append(result, n.Children...)
		}
	}
	return result
}

func filterChildren(node *models.AgentNode, allowed map[uint]bool) {
	var kept []*models.AgentNode
	for _, child := range node.Children {
		filterChildren(child, allowed)
		if allowed[child.AgentID] {
			kept = append(kept, child)
		} else {
			kept = append(kept, child.Children...)
		}
	}
	node.Children = kept
}

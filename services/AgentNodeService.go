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
	GetTrees(actor policies.AuthContext) ([]*agent_node.AgentNodeResponse, error)
}

type agentNodeService struct {
	repo        repositories.AgentNodeRepo
	scopeRepo   repositories.ScopeRepo
	agentPolicy policies.AgentPolicy
	scopePolicy policies.ScopePolicy
}

func NewAgentNodeService(
	repo repositories.AgentNodeRepo,
	scopeRepo repositories.ScopeRepo,
	agentPolicy policies.AgentPolicy,
	scopePolicy policies.ScopePolicy,
) AgentNodeService {
	return &agentNodeService{
		repo:        repo,
		scopeRepo:   scopeRepo,
		agentPolicy: agentPolicy,
		scopePolicy: scopePolicy,
	}
}

func (s *agentNodeService) CreateNode(request *agent_node.AgentNodeRequest, actor policies.AuthContext) (*agent_node.AgentNodeResponse, error) {
	if err := s.agentPolicy.CanCreateAgent(actor, request.ParentId); err != nil {
		return nil, err
	}

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

	if actor.Role == "OPERATOR" {
		if err := s.scopeRepo.AssignAgentToOperator(actor.UserID, newNode.AgentID); err != nil {
			return nil, err
		}
	}

	response := mappers.ToAgentNodeResponse(newNode)
	return &response, nil
}

func (s *agentNodeService) GetTrees(actor policies.AuthContext) ([]*agent_node.AgentNodeResponse, error) {
	treeModels, err := s.repo.GetTrees()
	if err != nil {
		return nil, err
	}

	scope, err := s.scopePolicy.AgentScope(actor)
	if err != nil {
		return nil, err
	}
	if scope.Unrestricted {
		return mappers.ToAgentNodePtrResponses(treeModels), nil
	}

	allowed := toSet(scope.IDs)
	pruned := filterRootsByAgentID(treeModels, allowed)
	return mappers.ToAgentNodePtrResponses(pruned), nil
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

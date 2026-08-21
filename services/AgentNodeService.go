package services

import (
	"errors"
	"example/go_backoffice/dto/agent_node"
	"example/go_backoffice/dto/user"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/models"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AgentNodeService interface {
	CreateNode(request *user.UserRequest, actor policies.AuthContext) (*agent_node.AgentNodeResponse, error)
	GetTree(actor policies.AuthContext) ([]*agent_node.AgentNodeResponse, error)
	DeleteNode(agentID uint, actor policies.AuthContext) error
}

type agentNodeService struct {
	db          *gorm.DB
	repo        repositories.AgentNodeRepo
	agentOpRepo repositories.AgentOperatorRepo
	scopeRepo   repositories.ScopeRepo
	userRepo    repositories.UserRepo
	policy      policies.UserPolicy
}

func NewAgentNodeService(
	db *gorm.DB,
	repo repositories.AgentNodeRepo,
	agentOpRepo repositories.AgentOperatorRepo,
	scopeRepo repositories.ScopeRepo,
	userRepo repositories.UserRepo,
	policy policies.UserPolicy,
) AgentNodeService {
	return &agentNodeService{
		db: db, repo: repo, agentOpRepo: agentOpRepo,
		scopeRepo: scopeRepo, userRepo: userRepo, policy: policy,
	}
}

func (s *agentNodeService) CreateNode(request *user.UserRequest, actor policies.AuthContext) (*agent_node.AgentNodeResponse, error) {
	if request == nil {
		return nil, errors.New("richiesta non valida")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	request.Password = string(hashed)

	var newNode *models.AgentNode

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		agentNodeRepo := s.repo.WithTx(tx)
		agentOpRepo := s.agentOpRepo.WithTx(tx)

		// Se è presente un ForeignId, il nuovo nodo va agganciato come figlio
		// del nodo dell'agente a cui il ForeignId fa riferimento
		var parentId *uint
		if request.ForeignId != nil {
			parentNode, err := agentNodeRepo.GetNodeByAgentID(*request.ForeignId)
			if err != nil {
				return err
			}
			parentId = &parentNode.ID
		}

		node := mappers.ToAgentNodeModel(request, parentId)

		if err := s.policy.Create(actor, node.Agent); err != nil {
			return err
		}

		if err := agentNodeRepo.Create(node); err != nil {
			return err
		}

		if actor.Role == enums.RoleOperator.String() {
			agentOp := &models.AgentOperator{
				OperatorID: actor.UserID,
				AgentID:    node.AgentID,
			}
			if _, err := agentOpRepo.AssignAgent(agentOp); err != nil {
				return err
			}
		}

		newNode = node
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	response := mappers.ToAgentNodeResponse(newNode)
	return &response, nil
}

func (s *agentNodeService) GetTree(actor policies.AuthContext) ([]*agent_node.AgentNodeResponse, error) {
	switch actor.Role {
	case enums.RoleAdmin.String():
		treeModels, err := s.repo.GetTrees()
		if err != nil {
			return nil, err
		}
		return mappers.ToAgentNodePtrResponses(treeModels), nil
	case enums.RoleOperator.String():
		treeModels, err := s.repo.GetFilteredTreeByOperator(actor.UserID)
		if err != nil {
			return nil, err
		}
		return mappers.ToAgentNodePtrResponses(treeModels), nil
	case enums.RoleAgent.String():
		treeModels, err := s.repo.GetFilteredTreeByAgent(actor.UserID)
		if err != nil {
			return nil, err
		}
		return mappers.ToAgentNodePtrResponses(treeModels), nil
	default:
		return nil, nil
	}
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

func (s *agentNodeService) DeleteNode(agentID uint, actor policies.AuthContext) error {
	target, err := s.userRepo.GetByIDAndRole(agentID, enums.RoleAgent.String())
	if err != nil {
		return err
	}

	if err := s.policy.Delete(actor, target); err != nil {
		return err
	}

	return s.repo.DeleteAgentNodeAndAgentByAgentID(agentID)
}

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
	GetTrees() ([]*agent_node.AgentNodeResponse, error)
}

type agentNodeService struct {
	db          *gorm.DB
	repo        repositories.AgentNodeRepo
	agentOpRepo repositories.AgentOperatorRepo
	scopeRepo   repositories.ScopeRepo
	policy      policies.UserPolicy
}

func NewAgentNodeService(
	db *gorm.DB,
	repo repositories.AgentNodeRepo,
	agentOpRepo repositories.AgentOperatorRepo,
	scopeRepo repositories.ScopeRepo,
	policy policies.UserPolicy,
) AgentNodeService {
	return &agentNodeService{
		db:          db,
		repo:        repo,
		agentOpRepo: agentOpRepo,
		scopeRepo:   scopeRepo,
		policy:      policy,
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

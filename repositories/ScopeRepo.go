package repositories

import (
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type ScopeRepo interface {
	AssignAgentToOperator(operatorID, agentID uint) error
	AssignAgencyToOperator(operatorID, agencyID uint) error
	IsAgentAssignedToOperator(operatorID, agentID uint) (bool, error)
	IsAgencyAssignedToOperator(operatorID, agencyID uint) (bool, error)
	AssignedAgentIDs(operatorID uint) ([]uint, error)
	AssignedAgencyIDs(operatorID uint) ([]uint, error)
	ChildrenAgentIds(agentId uint) ([]uint, error)
	GetAgentIDByNodeID(nodeID uint) (uint, error)
}

type scopeRepo struct {
	db *gorm.DB
}

func NewScopeRepository(db *gorm.DB) ScopeRepo {
	return &scopeRepo{db: db}
}

func (r *scopeRepo) AssignAgentToOperator(operatorID, agentID uint) error {
	return r.db.Create(&models.AgentOperator{OperatorID: operatorID, AgentID: agentID}).Error
}

func (r *scopeRepo) AssignAgencyToOperator(operatorID, agencyID uint) error {
	return r.db.Create(&models.AgencyOperator{OperatorID: operatorID, AgencyID: agencyID}).Error
}

func (r *scopeRepo) IsAgentAssignedToOperator(operatorID, agentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.AgentOperator{}).
		Where("operator_id = ? AND agent_id = ?", operatorID, agentID).
		Count(&count).Error
	return count > 0, err
}

func (r *scopeRepo) IsAgencyAssignedToOperator(operatorID, agencyID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.AgencyOperator{}).
		Where("operator_id = ? AND agency_id = ?", operatorID, agencyID).
		Count(&count).Error
	return count > 0, err
}

func (r *scopeRepo) AssignedAgentIDs(operatorID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&models.AgentOperator{}).
		Where("operator_id = ?", operatorID).
		Pluck("agent_id", &ids).Error
	return ids, err
}

func (r *scopeRepo) AssignedAgencyIDs(operatorID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&models.AgencyOperator{}).
		Where("operator_id = ?", operatorID).
		Pluck("agency_id", &ids).Error
	return ids, err
}

func (r *scopeRepo) ChildrenAgentIds(agentId uint) ([]uint, error) {
	// Risolve il nodo (AgentNode) associato all'AgentID di partenza
	var rootNode models.AgentNode
	if err := r.db.Where("agent_id = ?", agentId).First(&rootNode).Error; err != nil {
		return nil, err
	}

	// Includiamo l'agente di partenza
	ids := []uint{agentId}

	// Attraversamento ricorsivo dell'albero via AgentNode.ID (non via AgentID)
	var fetchDescendants func(currentNodeID uint) error
	fetchDescendants = func(currentNodeID uint) error {
		var children []models.AgentNode

		err := r.db.
			Where("parent_id = ?", currentNodeID).
			Find(&children).Error
		if err != nil {
			return err
		}

		for _, child := range children {
			ids = append(ids, child.AgentID)
			if err := fetchDescendants(child.ID); err != nil {
				return err
			}
		}

		return nil
	}

	if err := fetchDescendants(rootNode.ID); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *scopeRepo) GetAgentIDByNodeID(nodeID uint) (uint, error) {
	var node models.AgentNode
	if err := r.db.Select("agent_id").First(&node, nodeID).Error; err != nil {
		return 0, err
	}
	return node.AgentID, nil
}

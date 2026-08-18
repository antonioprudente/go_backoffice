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
	GetAgentIDByNodeID(nodeID uint) (uint, error)
	NodeChildrenAgentIds(agentID uint) ([]uint, error)
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

func (r *scopeRepo) GetAgentIDByNodeID(nodeID uint) (uint, error) {
	var node models.AgentNode
	if err := r.db.Select("agent_id").First(&node, nodeID).Error; err != nil {
		return 0, err
	}
	return node.AgentID, nil
}

func (r *scopeRepo) NodeChildrenAgentIds(agentID uint) ([]uint, error) {
	var node models.AgentNode
	if err := r.db.Where("agent_id = ?", agentID).First(&node).Error; err != nil {
		return nil, err
	}

	var children []uint
	err := r.db.Model(&models.AgentNode{}).
		Where("lft > ? AND rgt < ?", node.Lft, node.Rgt).
		Pluck("agent_id", &children).Error
	if err != nil {
		return nil, err
	}

	return children, nil
}

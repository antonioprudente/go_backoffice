package repositories

import (
	"example/go_backoffice/models"
	"fmt"

	"gorm.io/gorm"
)

type AgentNodeRepo interface {
}

type agentNodeRepo struct {
	db *gorm.DB
}

func NewAgentNodeRepository(db *gorm.DB) AgentNodeRepo {
	return &agentNodeRepo{db: db}
}

func (r *agentNodeRepo) Create(agentModel *models.AgentNode) error {
	return r.db.Create(agentModel).Error
}

func (r *agentNodeRepo) getParents(id *uint) ([]models.AgentNode, error) {
	if id == nil {
		return nil, fmt.Errorf("id non può essere nil")
	}

	var child models.AgentNode
	if err := r.db.First(&child, *id).Error; err != nil {
		return nil, err
	}

	var parents []models.AgentNode
	err := r.db.Where("lft < ? AND rgt > ?", child.Lft, child.Rgt).
		Order("lft ASC").
		Find(&parents).Error

	return parents, err
}

func (r *agentNodeRepo) getChildren(id *uint) ([]models.AgentNode, error) {
	if id == nil {
		return nil, fmt.Errorf("id non può essere nil")
	}

	var parent models.AgentNode
	if err := r.db.First(&parent, *id).Error; err != nil {
		return nil, err
	}

	var children []models.AgentNode
	err := r.db.Where("lft > ? AND rgt < ?", parent.Lft, parent.Rgt).
		Order("lft ASC").
		Find(&children).Error

	return children, err
}

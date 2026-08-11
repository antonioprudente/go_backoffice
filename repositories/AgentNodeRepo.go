package repositories

import (
	"example/go_backoffice/models"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AgentNodeRepo interface {
	Create(agentModel *models.AgentNode) error
	GetAncestors(id *uint) ([]models.AgentNode, error)
	GetDescendants(id *uint) ([]models.AgentNode, error)
}

type agentNodeRepo struct {
	db *gorm.DB
}

func NewAgentNodeRepository(db *gorm.DB) AgentNodeRepo {
	return &agentNodeRepo{db: db}
}

func (r *agentNodeRepo) Create(agentModel *models.AgentNode) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Crea prima l'utente collegato, se presente
		if agentModel.Agent != nil {
			if err := tx.Create(agentModel.Agent).Error; err != nil {
				return err
			}
			agentModel.AgentID = agentModel.Agent.ID
		}

		if agentModel.ParentID == nil {
			agentModel.Lft = 1
			agentModel.Rgt = 2
			return tx.Omit("Agent", "Parent").Create(agentModel).Error
		}

		var parent models.AgentNode
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&parent, *agentModel.ParentID).Error; err != nil {
			return err
		}

		parentRgt := parent.Rgt

		if err := tx.Model(&models.AgentNode{}).
			Where("rgt >= ?", parentRgt).
			Update("rgt", gorm.Expr("rgt + 2")).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.AgentNode{}).
			Where("lft >= ?", parentRgt).
			Update("lft", gorm.Expr("lft + 2")).Error; err != nil {
			return err
		}

		agentModel.Lft = parentRgt
		agentModel.Rgt = parentRgt + 1

		return tx.Omit("Agent", "Parent").Create(agentModel).Error
	})
}

func (r *agentNodeRepo) GetAncestors(id *uint) ([]models.AgentNode, error) {
	if id == nil {
		return nil, fmt.Errorf("id non può essere nil")
	}

	var child models.AgentNode
	if err := r.db.First(&child, *id).Error; err != nil {
		return nil, err
	}

	var ancestors []models.AgentNode
	err := r.db.Where("lft < ? AND rgt > ?", child.Lft, child.Rgt).
		Order("lft ASC").
		Find(&ancestors).Error

	return ancestors, err
}

func (r *agentNodeRepo) GetDescendants(id *uint) ([]models.AgentNode, error) {
	if id == nil {
		return nil, fmt.Errorf("id non può essere nil")
	}

	var parent models.AgentNode
	if err := r.db.First(&parent, *id).Error; err != nil {
		return nil, err
	}

	var descendants []models.AgentNode
	err := r.db.Where("lft > ? AND rgt < ?", parent.Lft, parent.Rgt).
		Order("lft ASC").
		Find(&descendants).Error

	return descendants, err
}

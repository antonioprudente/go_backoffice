package repositories

import (
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type AgentOperatorRepo interface {
	WithTx(tx *gorm.DB) AgentOperatorRepo

	AssignAgent(model *models.AgentOperator) (*models.AgentOperator, error)
	AssignAgentsMassive(agentOperator *[]models.AgentOperator) (*[]models.AgentOperator, error)
	GetByOperatorIDAndAgentID(operatorID uint, agentID uint) (models.AgentOperator, error)
	DeleteByAgentIDAndOperatorID(agentID uint, operatorID uint) (bool, error)
}

type agentOperatorRepo struct {
	db *gorm.DB
}

func NewAgentOperatorRepository(db *gorm.DB) AgentOperatorRepo {
	return &agentOperatorRepo{db: db}
}

func (r *agentOperatorRepo) WithTx(tx *gorm.DB) AgentOperatorRepo {
	return &agentOperatorRepo{db: tx}
}

func (r *agentOperatorRepo) AssignAgent(model *models.AgentOperator) (*models.AgentOperator, error) {
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *agentOperatorRepo) AssignAgentsMassive(agentOperator *[]models.AgentOperator) (*[]models.AgentOperator, error) {
	if len(*agentOperator) == 0 {
		return nil, nil
	}

	if err := r.db.Create(agentOperator).Error; err != nil {
		return nil, err
	}

	return agentOperator, nil
}

func (r *agentOperatorRepo) GetByOperatorIDAndAgentID(operatorID uint, agentID uint) (models.AgentOperator, error) {
	var agentOperator models.AgentOperator

	err := r.db.Where("operator_id = ? AND agent_id = ?", operatorID, agentID).First(&agentOperator).Error
	if err != nil {
		return models.AgentOperator{}, err
	}

	return agentOperator, nil
}

func (r *agentOperatorRepo) DeleteByAgentIDAndOperatorID(agentID uint, operatorID uint) (bool, error) {
	result := r.db.Where("operator_id = ? AND agent_id = ?", operatorID, agentID).Delete(&models.AgentOperator{})

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

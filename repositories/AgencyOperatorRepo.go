package repositories

import (
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type AgencyOperatorRepo interface {
	AssignAgency(model *models.AgencyOperator) (*models.AgencyOperator, error)
	AssignAgenciesMassive(agencyOperator []models.AgencyOperator) ([]models.AgencyOperator, error)
}

type agencyOperatorRepo struct {
	db *gorm.DB
}

func NewAgencyOperatorRepository(db *gorm.DB) AgencyOperatorRepo {
	return &agencyOperatorRepo{db: db}
}

func (r *agencyOperatorRepo) AssignAgency(model *models.AgencyOperator) (*models.AgencyOperator, error) {
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *agencyOperatorRepo) AssignAgenciesMassive(agencyOperator []models.AgencyOperator) ([]models.AgencyOperator, error) {
	if len(agencyOperator) == 0 {
		return nil, nil
	}

	if err := r.db.Create(&agencyOperator).Error; err != nil {
		return nil, err
	}

	return agencyOperator, nil
}

func (r *agentOperatorRepo) GetByOperatorIDAndAgencyID(operatorID uint, agencyID uint) (models.AgencyOperator, error) {
	var agencyOperator models.AgencyOperator

	err := r.db.Where("operator_id = ? AND agent_id = ?", operatorID, agencyID).First(&agencyOperator).Error
	if err != nil {
		return models.AgencyOperator{}, err
	}

	return agencyOperator, nil
}

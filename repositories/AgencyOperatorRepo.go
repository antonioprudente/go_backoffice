package repositories

import (
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type AgencyOperatorRepo interface {
	WithTx(tx *gorm.DB) AgencyOperatorRepo
	AssignAgency(model *models.AgencyOperator) (*models.AgencyOperator, error)
	AssignAgenciesMassive(agencyOperator *[]models.AgencyOperator) (*[]models.AgencyOperator, error)
	GetByOperatorIDAndAgencyID(operatorID uint, agencyID uint) (models.AgencyOperator, error)
	DeleteByAgencyIDAndOperatorID(agencyID uint, operatorID uint) (bool, error)
}

type agencyOperatorRepo struct {
	db *gorm.DB
}

func NewAgencyOperatorRepository(db *gorm.DB) AgencyOperatorRepo {
	return &agencyOperatorRepo{db: db}
}

func (r *agencyOperatorRepo) WithTx(tx *gorm.DB) AgencyOperatorRepo {
	return &agencyOperatorRepo{db: tx}
}

func (r *agencyOperatorRepo) AssignAgency(model *models.AgencyOperator) (*models.AgencyOperator, error) {
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *agencyOperatorRepo) AssignAgenciesMassive(agencyOperator *[]models.AgencyOperator) (*[]models.AgencyOperator, error) {
	if len(*agencyOperator) == 0 {
		return nil, nil
	}

	if err := r.db.Create(agencyOperator).Error; err != nil {
		return nil, err
	}

	return agencyOperator, nil
}

func (r *agencyOperatorRepo) GetByOperatorIDAndAgencyID(operatorID uint, agencyID uint) (models.AgencyOperator, error) {
	var agencyOperator models.AgencyOperator

	err := r.db.Where("operator_id = ? AND agent_id = ?", operatorID, agencyID).First(&agencyOperator).Error
	if err != nil {
		return models.AgencyOperator{}, err
	}

	return agencyOperator, nil
}

func (r *agencyOperatorRepo) DeleteByAgencyIDAndOperatorID(agencyID uint, operatorID uint) (bool, error) {
	result := r.db.Where("operator_id = ? AND agency_id = ?", operatorID, agencyID).Delete(&models.AgencyOperator{})

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

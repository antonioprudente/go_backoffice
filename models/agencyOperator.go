package models

import "gorm.io/gorm"

type AgencyOperator struct {
	gorm.Model
	AgencyID   uint `json:"agency_id" gorm:"primaryKey"`
	OperatorID uint `json:"operator_id" gorm:"primaryKey"`

	Agency   User `json:"agency,omitempty" gorm:"foreignKey:AgencyID"`
	Operator User `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

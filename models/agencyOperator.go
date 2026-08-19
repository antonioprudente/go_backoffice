package models

import "time"

type AgencyOperator struct {
	OperatorID uint      `gorm:"primaryKey;column:operator_id" json:"operator_id"`
	AgencyID   uint      `gorm:"primaryKey;column:agency_id" json:"agency_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (AgencyOperator) TableName() string {
	return "agency_operator"
}

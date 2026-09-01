package models

import "time"

type AgencyOperator struct {
	OperatorID uint      `gorm:"primaryKey;column:operator_id;autoIncrement:false;not null" json:"operator_id"`
	AgencyID   uint      `gorm:"primaryKey;column:agency_id;autoIncrement:false;not null" json:"agency_id"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;not null" json:"created_at"`
}

func (AgencyOperator) TableName() string {
	return "agency_operator"
}

package models

import "gorm.io/gorm"

type AgentOperator struct {
	gorm.Model
	AgentID    uint `json:"agent_id" gorm:"primaryKey"`
	OperatorID uint `json:"operator_id" gorm:"primaryKey"`

	Agent    *User `json:"agent,omitempty" gorm:"foreignKey:AgentID"`
	Operator *User `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

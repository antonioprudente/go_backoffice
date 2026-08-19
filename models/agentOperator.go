package models

import "time"

type AgentOperator struct {
	OperatorID uint      `gorm:"primaryKey;column:operator_id" json:"operator_id"`
	AgentID    uint      `gorm:"primaryKey;column:agent_id" json:"agent_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (AgentOperator) TableName() string {
	return "agent_operator"
}

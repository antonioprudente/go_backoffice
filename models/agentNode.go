package models

import "gorm.io/gorm"

type AgentNode struct {
	gorm.Model
	Lft       uint       `json:"lft" gorm:"not null"`
	Rgt       uint       `json:"rgt" gorm:"not null"`
	ParentId  uint       `json:"parent_id gorm:"not null"`
	AgentNode *AgentNode `json:"-" gorm:"foreignKey:ParentId"`
	AgentId   uint       `json:"agent_id" gorm:"not null"`
	Agent     *User      `json:"-" gorm:"foreignKey:AgentId"`
}

package models

import "gorm.io/gorm"

type AgentNode struct {
	gorm.Model
	Lft      uint  `json:"lft" gorm:"not null"`
	Rgt      uint  `json:"rgt" gorm:"not null"`
	ParentID *uint `json:"parent_id"`
	AgentID  uint  `json:"agent_id" gorm:"not null"`

	Parent *AgentNode `json:"-" gorm:"foreignKey:ParentID"`
	Agent  *User      `json:"-" gorm:"foreignKey:AgentID"`
}

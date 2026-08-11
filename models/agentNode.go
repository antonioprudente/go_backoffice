package models

import "gorm.io/gorm"

type AgentNode struct {
	gorm.Model
	Lft      uint  `json:"lft" gorm:"index; not null"`
	Rgt      uint  `json:"rgt" gorm:"index; not null"`
	ParentID *uint `json:"parent_id"`
	AgentID  uint  `json:"agent_id" gorm:"not null"`

	Parent *AgentNode `json:"-" gorm:"foreignKey:ParentID"`
	Agent  *User      `json:"-" gorm:"foreignKey:AgentID"`

	Children []*AgentNode `json:"children" gorm:"-"`
}

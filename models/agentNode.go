package models

import "gorm.io/gorm"

type AgentNode struct {
	gorm.Model
	Lft      uint  `json:"lft" gorm:"index;not null"`
	Rgt      uint  `json:"rgt" gorm:"index;not null"`
	ParentID *uint `json:"parent_id" gorm:"index"`
	AgentID  uint  `json:"agent_id" gorm:"index;not null"`

	Parent *AgentNode `json:"parent" gorm:"foreignKey:ParentID"`
	Agent  *User      `json:"agent" gorm:"foreignKey:AgentID"`

	Children []*AgentNode `json:"children" gorm:"-"`
	Agencies []*User      `json:"agencies" gorm:"-"`
}

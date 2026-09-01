package models

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	Content  string `gorm:"type:text" json:"content"`
	ActorID  uint   `gorm:"index;column:actor_id;not null" json:"actor_id"`
	TargetID uint   `gorm:"index;column:target_id;not null" json:"target_id"`

	// Relazioni verso la tabella users
	Actor  *User `gorm:"foreignKey:ActorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"actor,omitempty"`
	Target *User `gorm:"foreignKey:TargetID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"target,omitempty"`
}

func (Note) TableName() string {
	return "notes"
}

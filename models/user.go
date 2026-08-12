package models

import (
	"example/go_backoffice/enums"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string       `json:"first_name" gorm:"not null"`
	LastName  string       `json:"last_name" gorm:"not null"`
	Username  string       `json:"username" gorm:"unique;not null"`
	Role      enums.Role   `json:"role" gorm:"type:enum('ADMIN', 'OPERATOR', 'AGENT', 'AGENCY', 'USER')"`
	Status    enums.Status `json:"status" gorm:"type:enum('ACTIVE', 'SUSPENDED', 'BLOCKED', 'DEFAULT')"`
	Email     string       `json:"email" gorm:"unique;not null"`
	Password  string       `json:"password" gorm:"not null"`
	ForeignId *uint        `json:"foreign_id" gorm:"index"`
	Foreign   *User        `json:"foreign" gorm:"foreignKey:ForeignId"`
}

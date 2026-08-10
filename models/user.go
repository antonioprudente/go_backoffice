package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Username  string `json:"username" gorm:"unique;not null"`
	Status    string `json:"status" gorm:"type:enum('ACTIVE', 'SUSPENDED', 'BLOCKED', 'DEFAULT')"`
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"-" gorm:"not null"`
}

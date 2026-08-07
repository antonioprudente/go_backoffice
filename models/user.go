package models

import "time"

type User struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	FirstName string     `json:"first_name" gorm:"not null"`
	LastName  string     `json:"last_name" gorm:"not null"`
	Username  string     `json:"username" gorm:"unique;not null"`
	Email     string     `json:"email" gorm:"unique;not null"`
	Password  string     `json:"password" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

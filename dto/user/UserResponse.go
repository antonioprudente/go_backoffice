package user

import (
	"example/go_backoffice/enums"
	"time"
)

type UserResponse struct {
	Id          uint           `json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Username    string         `json:"username"`
	Role        enums.Role     `json:"role"`
	Status      enums.Status   `json:"status"`
	Email       string         `json:"email"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	ForeignID   *uint          `json:"foreign_id,omitempty"`
	Foreign     *UserResponse  `json:"foreign,omitempty"`
	LinkedUsers []UserResponse `json:"linked_users,omitempty"`
}

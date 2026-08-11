package user

import "example/go_backoffice/enums"

type UserResponse struct {
	Id        int           `json:"id"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Username  string        `json:"username"`
	Role      enums.Role    `json:"role"`
	Status    enums.Status  `json:"status"`
	Email     string        `json:"email"`
	ForeignId *uint         `json:"foreign_id"`
	Foreign   *UserResponse `json:"foreign"`
}

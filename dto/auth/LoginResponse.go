package auth

import "example/go_backoffice/dto/user"

type LoginResponse struct {
	Token string            `json:"token"`
	User  user.UserResponse `json:"user"`
}

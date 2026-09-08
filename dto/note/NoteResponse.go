package note

import (
	"example/go_backoffice/dto/user"
	"time"
)

type NoteResponse struct {
	ID             uint               `json:"id"`
	ActorUsername  *string            `json:"actor_username,omitempty"`
	TargetUsername *string            `json:"target_username,omitempty"`
	Actor          *user.UserResponse `json:"actor,omitempty"`
	Target         *user.UserResponse `json:"target,omitempty"`
	Content        string             `json:"content"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

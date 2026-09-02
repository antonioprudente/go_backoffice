package log

import (
	"example/go_backoffice/dto/user"
	"time"
)

type LogResponse struct {
	ID          uint               `json:"id"`
	ActorID     *uint              `json:"actor_id,omitempty"`
	TargetID    *uint              `json:"target_id,omitempty"`
	Actor       *user.UserResponse `json:"actor,omitempty"`
	Action      string             `json:"action"`
	Target      *user.UserResponse `json:"target,omitempty"`
	TargetType  string             `json:"target_type"`
	Description string             `json:"description"`
	Metadata    *string            `json:"metadata,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
}

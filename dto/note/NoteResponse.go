package note

import (
	"time"
)

type NoteResponse struct {
	ID        uint      `json:"id"`
	ActorID   uint      `json:"actor_id"`
	TargetID  uint      `json:"target_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

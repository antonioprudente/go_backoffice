package note

type NoteRequest struct {
	ActorID  uint   `json:"actor_id,omitempty"`
	TargetID uint   `json:"target_id"`
	Content  string `json:"content"`
}

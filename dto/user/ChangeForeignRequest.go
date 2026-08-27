package user

// dto/agent_node/AgentNodeRequest.go
type ChangeForeignRequest struct {
	UserID   uint  `json:"user_id"`
	TargetID *uint `json:"target_id,omitempty"`
}

package agent_node

// dto/agent_node/AgentNodeRequest.go
type MoveNodeRequest struct {
	AgentID  uint  `json:"agent_id"`
	TargetID *uint `json:"target_id,omitempty"`
}

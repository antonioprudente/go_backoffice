package agent_node

import "example/go_backoffice/dto/user"

// dto/agent_node/AgentNodeRequest.go
type AgentNodeRequest struct {
	ParentId *uint             `json:"parent_id,omitempty"`
	Agent    *user.UserRequest `json:"agent"`
}

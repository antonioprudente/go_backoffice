package agent_node

import "example/go_backoffice/dto/user"

type AgentNodeResponse struct {
	Id       uint                 `json:"id"`
	Lft      uint                 `json:"lft"`
	Rgt      uint                 `json:"rgt"`
	ParentId *uint                `json:"parent_id,omitempty"`
	Parent   *AgentNodeResponse   `json:"parent,omitempty"`
	Agent    user.UserResponse    `json:"agent"`
	Agencies []user.UserResponse  `json:"agencies,omitempty"`
	Children []*AgentNodeResponse `json:"children,omitempty"`
}

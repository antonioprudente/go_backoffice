// mappers/AgentNodeMapper.go
package mappers

import (
	"example/go_backoffice/dto/agent_node"
	"example/go_backoffice/models"
)

// ToAgentNodeModel converte una UserRequest nel model da persistere
func ToAgentNodeModel(req *agent_node.AgentNodeRequest) *models.AgentNode {
	return &models.AgentNode{
		ParentID: req.ParentId,
		Agent:    ToUserModel(req.Agent),
	}
}

// ToAgentNodeResponse converte un model nella response da esporre
func ToAgentNodeResponse(u *models.AgentNode) agent_node.AgentNodeResponse {
	resp := agent_node.AgentNodeResponse{
		Id:    u.ID,
		Lft:   u.Lft,
		Rgt:   u.Rgt,
		Agent: ToUserResponse(u.Agent),
	}

	if u.Parent != nil {
		parentResp := ToAgentNodeResponse(u.Parent)
		resp.Parent = &parentResp
	}

	return resp
}

// ToAgentNodeResponses converte uno slice di model, utile per GetAllNodes
func ToAgentNodeResponses(nodes []models.AgentNode) []agent_node.AgentNodeResponse {
	res := make([]agent_node.AgentNodeResponse, len(nodes))
	for i, n := range nodes {
		res[i] = ToAgentNodeResponse(&n)
	}
	return res
}

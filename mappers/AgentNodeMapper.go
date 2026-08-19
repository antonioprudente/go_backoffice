// mappers/AgentNodeMapper.go
package mappers

import (
	"example/go_backoffice/dto/agent_node"
	"example/go_backoffice/dto/user"
	"example/go_backoffice/models"
)

// ToAgentNodeModel converte una UserRequest nel model da persistere
func ToAgentNodeModel(req *user.UserRequest, parentId *uint) *models.AgentNode {
	return &models.AgentNode{
		ParentID: parentId,
		Agent:    ToUserModel(req),
	}
}

// ToAgentNodeResponse converte un model nella response da esporre
func ToAgentNodeResponse(u *models.AgentNode) agent_node.AgentNodeResponse {
	if u == nil {
		return agent_node.AgentNodeResponse{}
	}

	resp := agent_node.AgentNodeResponse{
		Id:       u.ID,
		Lft:      u.Lft,
		Rgt:      u.Rgt,
		ParentId: u.ParentID,
		Agent:    ToUserResponse(u.Agent),
	}

	// Mappatura delle agenzie associate all'agente (User con ForeignId == AgentID)
	resp.Agencies = make([]user.UserResponse, 0, len(u.Agencies))
	for _, agency := range u.Agencies {
		resp.Agencies = append(resp.Agencies, ToUserResponse(agency))
	}

	if u.Parent != nil {
		parentResp := ToAgentNodeResponse(u.Parent)
		resp.Parent = &parentResp
	}

	// Mappatura ricorsiva dei nodi figli (Children)
	if len(u.Children) > 0 {
		resp.Children = make([]*agent_node.AgentNodeResponse, len(u.Children))

		for i, child := range u.Children {
			childRes := ToAgentNodeResponse(child)
			resp.Children[i] = &childRes
		}
	}

	return resp
}

// ToAgentNodeResponses converte uno slice di model, utile per GetAllNodes
func ToAgentNodeResponses(nodes []*models.AgentNode) []agent_node.AgentNodeResponse {
	res := make([]agent_node.AgentNodeResponse, len(nodes))
	for i, n := range nodes {
		res[i] = ToAgentNodeResponse(n)
	}

	return res
}

func ToAgentNodePtrResponses(nodes []*models.AgentNode) []*agent_node.AgentNodeResponse {
	if nodes == nil {
		return nil
	}

	res := make([]*agent_node.AgentNodeResponse, len(nodes))
	for i, n := range nodes {
		// Mappiamo direttamente il puntatore passandolo a ToAgentNodeResponse
		nodeRes := ToAgentNodeResponse(n)
		res[i] = &nodeRes
	}

	return res
}

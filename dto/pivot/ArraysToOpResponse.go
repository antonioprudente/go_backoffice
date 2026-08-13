package pivot

import "example/go_backoffice/dto/user"

type ArraysToOpResponse struct {
	OperatorId user.UserResponse   `json:"operator_id"`
	AgentIds   []user.UserResponse `json:"agent_ids,omitempty"`
	AgencyIds  []user.UserResponse `json:"agency_ids,omitempty"`
}

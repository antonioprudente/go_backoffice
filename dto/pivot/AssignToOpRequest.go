package pivot

type AssignToOpRequest struct {
	OperatorId uint  `json:"operator_id"`
	AgentId    *uint `json:"agent_id,omitempty"`
	AgencyId   *uint `json:"agency_id,omitempty"`
}

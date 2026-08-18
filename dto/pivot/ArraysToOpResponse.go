package pivot

type ArraysToOpResponse struct {
	OperatorId uint    `json:"operator_id"`
	AgentIds   *[]uint `json:"agent_ids,omitempty"`
	AgencyIds  *[]uint `json:"agency_ids,omitempty"`
}

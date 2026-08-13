package pivot

type ArraysToOpRequest struct {
	OperatorId uint    `json:"operator_id"`
	AgentIds   *[]uint `json:"agent_ids,omitempty"`
	AgencyIds  *[]uint `json:"agency_ids,omitempty"`
}

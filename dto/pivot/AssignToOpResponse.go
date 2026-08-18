package pivot

type AssignToOpResponse struct {
	OperatorID uint  `json:"operator"`
	AgentID    *uint `json:"agent,omitempty"`
	AgencyID   *uint `json:"agency,omitempty"`
}

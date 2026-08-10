package models

// AgentOperator rappresenta l'associazione N:N tra Agent e Operator.
// Entrambi i campi fanno riferimento a User.ID (rispettivamente con Role=AGENT e Role=OPERATOR).
type AgentOperator struct {
	AgentID    uint `json:"agent_id" gorm:"primaryKey"`
	OperatorID uint `json:"operator_id" gorm:"primaryKey"`

	Agent    User `json:"agent,omitempty" gorm:"foreignKey:AgentID"`
	Operator User `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

// TableName forza il nome tabella, altrimenti GORM genererebbe "agent_operators"
// dal nome dello struct — qui è già corretto, ma lo esplicito per chiarezza.
func (AgentOperator) TableName() string {
	return "agents_operators"
}

package policies

import (
	"errors"

	"example/go_backoffice/enums"
	"example/go_backoffice/repositories"
)

var ErrForbidden = errors.New("operazione non consentita per il tuo ruolo")

type AgentPolicy interface {
	// CanCreateAgent verifica se actor può creare un nuovo AgentNode sotto parentNodeID
	CanCreateAgent(actor AuthContext, parentNodeID *uint) error
}

type agentPolicy struct {
	agentNodeRepo repositories.AgentNodeRepo
}

func NewAgentPolicy(agentNodeRepo repositories.AgentNodeRepo) AgentPolicy {
	return &agentPolicy{agentNodeRepo: agentNodeRepo}
}

func (p *agentPolicy) CanCreateAgent(actor AuthContext, parentNodeID *uint) error {
	switch enums.Role(actor.Role) {
	case enums.RoleAdmin:
		return nil

	case enums.RoleOperator:
		// l'operatore può sempre creare un nuovo agente: verrà assegnato a lui
		return nil

	case enums.RoleAgent:
		if parentNodeID == nil {
			return errors.New("un agente non può creare un nodo radice")
		}
		allowed, err := p.agentNodeRepo.IsDescendantOrSelf(actor.UserID, *parentNodeID)
		if err != nil {
			return err
		}
		if !allowed {
			return errors.New("non puoi creare agenti fuori dalla tua gerarchia")
		}
		return nil

	default:
		return ErrForbidden
	}
}

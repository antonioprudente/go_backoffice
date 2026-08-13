package policies

import (
	"errors"

	"example/go_backoffice/enums"
	"example/go_backoffice/repositories"
)

type AgencyPolicy interface {
	// CanCreateAgency verifica se actor può creare un'AGENCY collegata a foreignAgentID
	CanCreateAgency(actor AuthContext, foreignAgentID *uint) error
}

type agencyPolicy struct {
	agentNodeRepo repositories.AgentNodeRepo
}

func NewAgencyPolicy(agentNodeRepo repositories.AgentNodeRepo) AgencyPolicy {
	return &agencyPolicy{agentNodeRepo: agentNodeRepo}
}

func (p *agencyPolicy) CanCreateAgency(actor AuthContext, foreignAgentID *uint) error {
	switch enums.Role(actor.Role) {
	case enums.RoleAdmin:
		return nil

	case enums.RoleOperator:
		return nil

	case enums.RoleAgent:
		if foreignAgentID == nil {
			return errors.New("specificare l'agente a cui collegare l'agenzia")
		}
		allowed, err := p.agentNodeRepo.IsAgentDescendantOrSelf(actor.UserID, *foreignAgentID)
		if err != nil {
			return err
		}
		if !allowed {
			return errors.New("non puoi creare agenzie fuori dalla tua gerarchia")
		}
		return nil

	default:
		return ErrForbidden
	}
}

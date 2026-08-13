package policies

import (
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
)

type UserPolicy interface {
	// CanViewUser verifica se actor può visualizzare la scheda dell'utente target,
	// qualunque sia il ruolo di actor e di target.
	CanViewUser(actor AuthContext, target *models.User) error
}

type userPolicy struct {
	agentNodeRepo repositories.AgentNodeRepo
	scopeRepo     repositories.ScopeRepo
	userRepo      repositories.UserRepo
}

func NewUserPolicy(
	agentNodeRepo repositories.AgentNodeRepo,
	scopeRepo repositories.ScopeRepo,
	userRepo repositories.UserRepo,
) UserPolicy {
	return &userPolicy{
		agentNodeRepo: agentNodeRepo,
		scopeRepo:     scopeRepo,
		userRepo:      userRepo,
	}
}

// Regole:
//   - ADMIN: tutti
//   - OPERATOR: se stesso, gli AGENT/AGENCY a lui assegnati, e gli USER
//     appartenenti alle AGENCY a lui assegnate
//   - AGENT: se stesso e i propri discendenti (AGENT), le AGENCY collegate a sé
//     o a un discendente, e gli USER di quelle AGENCY
//   - AGENCY: se stessa e i propri USER (target.ForeignId == actor.UserID)
func (p *userPolicy) CanViewUser(actor AuthContext, target *models.User) error {
	if target == nil {
		return ErrForbidden
	}

	switch enums.Role(actor.Role) {
	case enums.RoleAdmin:
		return nil

	case enums.RoleOperator:
		return p.operatorCanView(actor, target)

	case enums.RoleAgent:
		return p.agentCanView(actor, target)

	case enums.RoleAgency:
		return p.agencyCanView(actor, target)

	default:
		return ErrForbidden
	}
}

func (p *userPolicy) operatorCanView(actor AuthContext, target *models.User) error {
	if target.Role == enums.RoleOperator && target.ID == actor.UserID {
		return nil
	}

	switch target.Role {
	case enums.RoleAgent:
		assigned, err := p.scopeRepo.IsAgentAssignedToOperator(actor.UserID, target.ID)
		if err != nil {
			return err
		}
		if assigned {
			return nil
		}

	case enums.RoleAgency:
		assigned, err := p.scopeRepo.IsAgencyAssignedToOperator(actor.UserID, target.ID)
		if err != nil {
			return err
		}
		if assigned {
			return nil
		}

	case enums.RoleUser:
		if target.ForeignId == nil {
			return ErrForbidden
		}
		assigned, err := p.scopeRepo.IsAgencyAssignedToOperator(actor.UserID, *target.ForeignId)
		if err != nil {
			return err
		}
		if assigned {
			return nil
		}
	}

	return ErrForbidden
}

func (p *userPolicy) agentCanView(actor AuthContext, target *models.User) error {
	switch target.Role {
	case enums.RoleAgent:
		allowed, err := p.agentNodeRepo.IsAgentDescendantOrSelf(actor.UserID, target.ID)
		if err != nil {
			return err
		}
		if allowed {
			return nil
		}

	case enums.RoleAgency:
		if target.ForeignId == nil {
			return ErrForbidden
		}
		allowed, err := p.agentNodeRepo.IsAgentDescendantOrSelf(actor.UserID, *target.ForeignId)
		if err != nil {
			return err
		}
		if allowed {
			return nil
		}

	case enums.RoleUser:
		if target.ForeignId == nil {
			return ErrForbidden
		}
		agency, err := p.userRepo.GetByIDAndRole(*target.ForeignId, enums.RoleAgency.String())
		if err != nil {
			return err
		}
		if agency.ForeignId == nil {
			return ErrForbidden
		}
		allowed, err := p.agentNodeRepo.IsAgentDescendantOrSelf(actor.UserID, *agency.ForeignId)
		if err != nil {
			return err
		}
		if allowed {
			return nil
		}
	}

	return ErrForbidden
}

func (p *userPolicy) agencyCanView(actor AuthContext, target *models.User) error {
	if target.Role == enums.RoleAgency && target.ID == actor.UserID {
		return nil
	}
	if target.Role == enums.RoleUser && target.ForeignId != nil && *target.ForeignId == actor.UserID {
		return nil
	}
	return ErrForbidden
}

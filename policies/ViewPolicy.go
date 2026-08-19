package policies

import (
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
	"slices"
)

type ViewPolicy struct {
	scopeRepo repositories.ScopeRepo
	userRepo  repositories.UserRepo
}

func NewViewPolicy(scopeRepo repositories.ScopeRepo, userRepo repositories.UserRepo) *ViewPolicy {
	return &ViewPolicy{scopeRepo: scopeRepo, userRepo: userRepo}
}

// Check ritorna nil se actor può vedere target, ErrForbidden se il permesso
// è negato (denial legittimo di business), oppure un altro errore per
// problemi tecnici (query fallita, ruolo non gestito, dato mancante).
func (p *ViewPolicy) Check(actor AuthContext, target *models.User) error {
	if actor.UserID == target.ID {
		return nil
	}
	switch actor.Role {
	case enums.RoleAdmin.String():
		return nil

	case enums.RoleOperator.String():
		return p.viewAsOperator(actor, target)

	case enums.RoleAgent.String():
		return p.viewAsAgent(actor, target)

	case enums.RoleAgency.String():
		if target.Role == enums.RoleUser && belongsTo(target, actor.UserID) {
			return nil
		}
		return ErrForbidden
	}
	return ErrUnknownRole
}

func (p *ViewPolicy) viewAsOperator(actor AuthContext, target *models.User) error {
	switch target.Role {
	case enums.RoleAgent:
		assigned, err := p.scopeRepo.IsAgentAssignedToOperator(actor.UserID, target.ID)
		if err != nil {
			return err
		}
		if !assigned {
			return ErrForbidden
		}
		return nil

	case enums.RoleAgency:
		assigned, err := p.scopeRepo.IsAgencyAssignedToOperator(actor.UserID, target.ID)
		if err != nil {
			return err
		}
		if !assigned {
			return ErrForbidden
		}
		return nil

	case enums.RoleUser:
		if target.ForeignId == nil {
			return ErrMissingRelation
		}
		assigned, err := p.scopeRepo.IsAgencyAssignedToOperator(actor.UserID, *target.ForeignId)
		if err != nil {
			return err
		}
		if !assigned {
			return ErrForbidden
		}
		return nil
	}
	return ErrUnknownRole
}

func (p *ViewPolicy) viewAsAgent(actor AuthContext, target *models.User) error {
	switch target.Role {
	case enums.RoleAgent, enums.RoleAgency:
		if target.ForeignId == nil {
			return ErrForbidden
		}
		children, err := p.scopeRepo.NodeChildrenAgentIds(actor.UserID)
		if err != nil {
			return err
		}
		if actor.UserID != *target.ForeignId && !slices.Contains(children, *target.ForeignId) {
			return ErrForbidden
		}
		return nil

	case enums.RoleUser:
		agency, err := p.userRepo.GetByIDAndRole(*target.ForeignId, enums.RoleAgency.String())
		if err != nil {
			return err
		}

		children, err := p.scopeRepo.NodeChildrenAgentIds(*agency.ForeignId)
		if err != nil {
			return err
		}

		if !slices.Contains(children, *target.ForeignId) {
			return ErrForbidden
		}
		return nil
	}

	return ErrUnknownRole
}

package policies

import (
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
)

type UpdateStatusPolicy struct {
	scopeRepo repositories.ScopeRepo
}

func NewUpdateStatusPolicy(scopeRepo repositories.ScopeRepo) *UpdateStatusPolicy {
	return &UpdateStatusPolicy{scopeRepo: scopeRepo}
}

// Check ritorna nil se actor può modificare target, ErrForbidden altrimenti.
func (p *UpdateStatusPolicy) Check(actor AuthContext, target *models.User) error {
	if actor.UserID == target.ID {
		return ErrForbidden
	}

	switch actor.Role {
	case enums.RoleAdmin.String():
		return nil

	case enums.RoleOperator.String():
		return p.updateStatusAsOperator(actor, target)

	case enums.RoleAgent.String():
		return p.updateStatusAsAgent(actor, target)

	case enums.RoleAgency.String():
		return nil
	}

	return ErrUnknownRole
}

func (p *UpdateStatusPolicy) updateStatusAsOperator(actor AuthContext, target *models.User) error {
	switch target.Role {
	case enums.RoleOperator:
		return ErrForbidden

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

func (p *UpdateStatusPolicy) updateStatusAsAgent(actor AuthContext, target *models.User) error {
	switch target.Role {
	case enums.RoleOperator:
		return ErrForbidden

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

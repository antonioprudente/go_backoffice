package policies

import (
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
	"slices"
)

type UpdatePolicy struct {
	scopeRepo repositories.ScopeRepo
	userRepo  repositories.UserRepo
}

func NewUpdatePolicy(scopeRepo repositories.ScopeRepo, userRepo repositories.UserRepo) *UpdatePolicy {
	return &UpdatePolicy{scopeRepo: scopeRepo, userRepo: userRepo}
}

// Check ritorna nil se actor può modificare target, ErrForbidden altrimenti.
func (p *UpdatePolicy) Check(actor AuthContext, target *models.User) error {
	if actor.UserID == target.ID {
		return nil
	}

	switch actor.Role {
	case enums.RoleAdmin.String():
		return nil

	case enums.RoleOperator.String():
		return p.updateAsOperator(actor, target)

	case enums.RoleAgent.String():
		return p.updateAsAgent(actor, target)
	}

	return ErrForbidden
}

func (p *UpdatePolicy) updateAsOperator(actor AuthContext, target *models.User) error {
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
	return ErrNotImplemented
}

func (p *UpdatePolicy) updateAsAgent(actor AuthContext, target *models.User) error {
	switch target.Role {
	case enums.RoleOperator:
		return ErrForbidden

	case enums.RoleAgent, enums.RoleAgency:
		if target.ForeignId == nil {
			if target.Role == enums.RoleAgent {
				return ErrForbidden
			}
			if target.Role == enums.RoleAgency {
				return ErrMissingRelation
			}
		}
		children, err := p.scopeRepo.NodeChildrenAndSelfAgentIds(actor.UserID)
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
		children, err := p.scopeRepo.NodeChildrenAndSelfAgentIds(actor.UserID)
		if err != nil {
			return err
		}
		if !slices.Contains(children, *agency.ForeignId) {
			return ErrForbidden
		}
		return nil
	}

	return ErrUnknownRole
}

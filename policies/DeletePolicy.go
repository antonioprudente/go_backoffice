package policies

import (
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
	"slices"
)

type DeletePolicy struct {
	scopeRepo repositories.ScopeRepo
	userRepo  repositories.UserRepo
}

func NewDeletePolicy(scopeRepo repositories.ScopeRepo, userRepo repositories.UserRepo) *DeletePolicy {
	return &DeletePolicy{scopeRepo: scopeRepo, userRepo: userRepo}
}

// Check ritorna nil se actor può eliminare target, ErrForbidden altrimenti.
func (p *DeletePolicy) Check(actor AuthContext, target *models.User) error {
	if actor.UserID == target.ID {
		return ErrForbidden
	}
	switch actor.Role {
	case enums.RoleAdmin.String():
		return nil

	case enums.RoleOperator.String():
		return p.deleteAsOperator(actor, target)

	case enums.RoleAgent.String():
		return p.deleteAsAgent(actor, target)
	}

	return ErrForbidden
}

func (p *DeletePolicy) deleteAsOperator(actor AuthContext, target *models.User) error {
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

func (p *DeletePolicy) deleteAsAgent(actor AuthContext, target *models.User) error {
	switch target.Role {
	case enums.RoleOperator:
		return ErrForbidden

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

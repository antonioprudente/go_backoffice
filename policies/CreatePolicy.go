package policies

import (
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
	"slices"
)

type CreatePolicy struct {
	scopeRepo repositories.ScopeRepo
	userRepo  repositories.UserRepo
}

func NewCreatePolicy(scopeRepo repositories.ScopeRepo, userRepo repositories.UserRepo) *CreatePolicy {
	return &CreatePolicy{scopeRepo: scopeRepo, userRepo: userRepo}
}

// Check ritorna nil se actor può creare un utente con i dati di target,
// ErrForbidden altrimenti.
func (p *CreatePolicy) Check(actor AuthContext, target *models.User) error {
	switch actor.Role {
	case enums.RoleAdmin.String():
		return nil

	case enums.RoleOperator.String():
		return p.createAsOperator(actor, target)

	case enums.RoleAgent.String():
		return p.createAsAgent(actor, target)
	}

	return ErrForbidden
}

func (p *CreatePolicy) createAsOperator(actor AuthContext, target *models.User) error {
	switch target.Role {
	case enums.RoleOperator:
		return ErrForbidden

	case enums.RoleAgent:
		return nil

	case enums.RoleAgency:
		assigned, err := p.scopeRepo.IsAgentAssignedToOperator(actor.UserID, *target.ForeignId)
		if err != nil {
			return err
		}
		if !assigned {
			return ErrForbidden
		}
		return nil

	case enums.RoleUser:
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

func (p *CreatePolicy) createAsAgent(actor AuthContext, target *models.User) error {
	switch target.Role {
	case enums.RoleOperator:
		return ErrForbidden

	case enums.RoleAgent, enums.RoleAgency:
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

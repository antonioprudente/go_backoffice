package policies

import (
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
)

type UpdatePolicy struct {
	scopeRepo repositories.ScopeRepo
}

func NewUpdatePolicy(scopeRepo repositories.ScopeRepo) *UpdatePolicy {
	return &UpdatePolicy{scopeRepo: scopeRepo}
}

// Check ritorna nil se actor può modificare target, ErrForbidden altrimenti.
func (p *UpdatePolicy) Check(actor AuthContext, target *models.User) error {
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
	// TODO: logica di business per OPERATOR
	return ErrNotImplemented
}

func (p *UpdatePolicy) updateAsAgent(actor AuthContext, target *models.User) error {
	// TODO: logica di business per AGENT
	return ErrNotImplemented
}

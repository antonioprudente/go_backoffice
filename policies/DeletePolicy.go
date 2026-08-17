package policies

import (
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
)

type DeletePolicy struct {
	scopeRepo repositories.ScopeRepo
}

func NewDeletePolicy(scopeRepo repositories.ScopeRepo) *DeletePolicy {
	return &DeletePolicy{scopeRepo: scopeRepo}
}

// Check ritorna nil se actor può eliminare target, ErrForbidden altrimenti.
func (p *DeletePolicy) Check(actor AuthContext, target *models.User) error {
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
	// TODO: logica di business per OPERATOR
	return ErrNotImplemented
}

func (p *DeletePolicy) deleteAsAgent(actor AuthContext, target *models.User) error {
	// TODO: logica di business per AGENT
	return ErrNotImplemented
}

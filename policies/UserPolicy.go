package policies

import (
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
)

type UserPolicy interface {
	View(actor AuthContext, target *models.User) error
	Create(actor AuthContext, target *models.User) error
	Update(actor AuthContext, target *models.User) error
	UpdateStatus(actor AuthContext, target *models.User) error
	Delete(actor AuthContext, target *models.User) error
}

type userPolicy struct {
	viewPolicy         *ViewPolicy
	createPolicy       *CreatePolicy
	updatePolicy       *UpdatePolicy
	updateStatusPolicy *UpdateStatusPolicy
	deletePolicy       *DeletePolicy
}

func NewUserPolicy(scopeRepo repositories.ScopeRepo, userRepo repositories.UserRepo) UserPolicy {
	return &userPolicy{
		viewPolicy:         NewViewPolicy(scopeRepo, userRepo),
		createPolicy:       NewCreatePolicy(scopeRepo, userRepo),
		updatePolicy:       NewUpdatePolicy(scopeRepo),
		updateStatusPolicy: NewUpdateStatusPolicy(scopeRepo),
		deletePolicy:       NewDeletePolicy(scopeRepo),
	}
}

func (p *userPolicy) View(actor AuthContext, target *models.User) error {
	return p.viewPolicy.Check(actor, target)
}

func (p *userPolicy) Create(actor AuthContext, target *models.User) error {
	return p.createPolicy.Check(actor, target)
}

func (p *userPolicy) Update(actor AuthContext, target *models.User) error {
	return p.updatePolicy.Check(actor, target)
}

func (p *userPolicy) UpdateStatus(actor AuthContext, target *models.User) error {
	return p.updateStatusPolicy.Check(actor, target)
}

func (p *userPolicy) Delete(actor AuthContext, target *models.User) error {
	return p.deletePolicy.Check(actor, target)
}

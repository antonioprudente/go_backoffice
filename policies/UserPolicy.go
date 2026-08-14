package policies

import (
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
)

type UserPolicy interface {
	View(actor AuthContext, target models.User) bool
}

type userPolicy struct {
	userRepo     repositories.UserRepo
	agentOpRepo  repositories.AgentOperatorRepo
	agencyOpRepo repositories.AgencyOperatorRepo
	scopeRepo    repositories.ScopeRepo
}

func NewUserPolicy(
	userRepo repositories.UserRepo,
	agentOpRepo repositories.AgentOperatorRepo,
	agencyOpRepo repositories.AgencyOperatorRepo,
	scopeRepo repositories.ScopeRepo,
) UserPolicy {
	return &userPolicy{
		userRepo:     userRepo,
		agentOpRepo:  agentOpRepo,
		agencyOpRepo: agencyOpRepo,
		scopeRepo:    scopeRepo,
	}
}

func (p *userPolicy) View(actor AuthContext, target models.User) bool {
	if actor.UserID == target.ID {
		return true
	}

	switch actor.Role {
	case enums.RoleAdmin.String():
		return true

	case enums.RoleOperator.String():
		switch target.Role {
		case enums.RoleAgent:
			isAssigned, _ := p.scopeRepo.IsAgentAssignedToOperator(actor.UserID, target.ID)
			return isAssigned

		case enums.RoleAgency:
			isAssigned, _ := p.scopeRepo.IsAgencyAssignedToOperator(actor.UserID, target.ID)
			return isAssigned

		case enums.RoleUser:
			isAssigned, _ := p.scopeRepo.IsAgencyAssignedToOperator(actor.UserID, target.Foreign.ID)
			return isAssigned
		}
		return false
	case enums.RoleAgent.String():
		switch target.Role {
		case enums.RoleAgency:
			return target.ForeignId == &actor.UserID

		case enums.RoleUser:
			return target.Foreign.ForeignId == &actor.UserID
		}
		return false
	}
	return false
}

func Create(actor AuthContext) bool {
	return false
}
func Update()
func Delete()

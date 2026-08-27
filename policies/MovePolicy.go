package policies

import (
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
	"slices"
)

type MovePolicy struct {
	scopeRepo repositories.ScopeRepo
	userRepo  repositories.UserRepo
}

func NewMovePolicy(scopeRepo repositories.ScopeRepo, userRepo repositories.UserRepo) *MovePolicy {
	return &MovePolicy{scopeRepo: scopeRepo, userRepo: userRepo}
}

// Check ritorna nil se actor può spostare il nodo "target" sotto "newParent",
// ErrForbidden altrimenti. newParent == nil significa "porta il nodo a root":
// operazione consentita solo all'ADMIN.
func (p *MovePolicy) Check(actor AuthContext, target *models.User, newParent *models.User) error {
	if target == nil {
		return ErrMissingRelation
	}
	if actor.UserID == target.ID {
		return ErrForbidden // non ha senso spostare se stessi
	}

	switch actor.Role {
	case enums.RoleAdmin.String():
		return nil

	case enums.RoleOperator.String():
		if newParent == nil {
			return ErrForbidden // l'operatore non può portare un nodo a root
		}
		return p.moveAsOperator(actor, target, newParent)

	case enums.RoleAgent.String():
		if newParent == nil {
			return ErrForbidden // idem per l'agente
		}
		return p.moveAsAgent(actor, target, newParent)
	}

	return ErrForbidden
}

// moveAsOperator: sia il nodo spostato che quello di destinazione devono
// essere agenti assegnati all'operatore.
func (p *MovePolicy) moveAsOperator(actor AuthContext, target *models.User, newParent *models.User) error {
	nodeAssigned, err := p.scopeRepo.IsAgentAssignedToOperator(actor.UserID, target.ID)
	if err != nil {
		return err
	}
	if !nodeAssigned {
		return ErrForbidden
	}

	destAssigned, err := p.scopeRepo.IsAgentAssignedToOperator(actor.UserID, newParent.ID)
	if err != nil {
		return err
	}
	if !destAssigned {
		return ErrForbidden
	}

	return nil
}

// moveAsAgent: sia il nodo spostato che la destinazione devono trovarsi
// nel sottoalbero dell'agente (sé stesso incluso).
func (p *MovePolicy) moveAsAgent(actor AuthContext, target *models.User, newParent *models.User) error {
	// Vieta di spostare l'agente loggato stesso
	if target.ID == actor.UserID {
		return ErrForbidden
	}

	descendants, err := p.scopeRepo.NodeChildrenAndSelfAgentIds(actor.UserID)
	if err != nil {
		return err
	}

	if !slices.Contains(descendants, target.ID) {
		return ErrForbidden
	}

	if !slices.Contains(descendants, newParent.ID) {
		return ErrForbidden
	}

	if actor.UserID == target.ID {
		return ErrForbidden
	}

	return nil
}

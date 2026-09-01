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
	switch target.Role {
	case enums.RoleAgency:
		if newParent.Role != enums.RoleAgent {
			return ErrForbidden
		}

		trgtAgencyAssigned, err := p.scopeRepo.IsAgencyAssignedToOperator(actor.UserID, target.ID)
		if err != nil {
			return err
		}

		if !trgtAgencyAssigned {
			return ErrForbidden
		}

		newPrntAssigned, err := p.scopeRepo.IsAgentAssignedToOperator(actor.UserID, newParent.ID)
		if err != nil {
			return err
		}

		if !newPrntAssigned {
			return ErrForbidden
		}
		return nil

	case enums.RoleAgent:
		if newParent.Role != enums.RoleAgent {
			return ErrForbidden
		}

		targetAssigned, err := p.scopeRepo.IsAgentAssignedToOperator(actor.UserID, target.ID)
		if err != nil {
			return err
		}

		if !targetAssigned {
			return ErrForbidden
		}

		newParentAssigned, err := p.scopeRepo.IsAgentAssignedToOperator(actor.UserID, newParent.ID)
		if err != nil {
			return err
		}

		if !newParentAssigned {
			return ErrForbidden
		}
		return nil

	case enums.RoleUser:
		if newParent.Role != enums.RoleAgency {
			return ErrForbidden
		}
		targetAssigned, err := p.scopeRepo.IsAgencyAssignedToOperator(actor.UserID, *target.ForeignID)

		if err != nil {
			return err
		}
		if !targetAssigned {
			return ErrForbidden
		}

		newParentAssigned, err := p.scopeRepo.IsAgencyAssignedToOperator(actor.UserID, target.ID)
		if !newParentAssigned {
			return ErrForbidden
		}
		return nil
	}
	return ErrUnknownRole
}

// moveAsAgent: sia il nodo spostato che la destinazione devono trovarsi
// nel sottoalbero dell'agente (sé stesso incluso).
func (p *MovePolicy) moveAsAgent(actor AuthContext, target *models.User, newParent *models.User) error {
	descendants, err := p.scopeRepo.NodeChildrenAndSelfAgentIds(actor.UserID)
	if err != nil {
		return err
	}

	switch target.Role {
	case enums.RoleAgency:
		if target.ForeignID == nil {
			return ErrMissingRelation
		}

		// L'agenzia "appartiene" al sottoalbero tramite l'agente a cui è
		// agganciata (ForeignID), non tramite il proprio ID: sia l'agente
		// proprietario attuale sia il nuovo agente di destinazione devono
		// trovarsi nel tuo sottoalbero.
		if !slices.Contains(descendants, *target.ForeignID) {
			return ErrForbidden
		}

		if newParent == nil {
			return ErrForbidden
		}

		if !slices.Contains(descendants, newParent.ID) {
			return ErrForbidden
		}
		return nil

	case enums.RoleAgent:
		if target.ID == actor.UserID {
			return ErrForbidden
		}

		if newParent == nil {
			return ErrForbidden
		}

		if !slices.Contains(descendants, newParent.ID) || !slices.Contains(descendants, target.ID) {
			return ErrForbidden
		}
		return nil

	case enums.RoleUser:
		if newParent == nil {
			return ErrMissingRelation
		}

		nowParent, err := p.userRepo.GetByIDAndRole(*target.ForeignID, enums.RoleUser.String())
		if err != nil {
			return err
		}

		if !slices.Contains(descendants, *nowParent.ForeignID) || !slices.Contains(descendants, newParent.ID) {
			return ErrForbidden
		}
	}

	return ErrUnknownRole
}

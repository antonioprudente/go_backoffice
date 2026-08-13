package policies

import (
	"errors"

	"example/go_backoffice/enums"
	"example/go_backoffice/repositories"
)

// Scope descrive come un service deve filtrare una lista di risorse
// in base a chi sta chiedendo.
type Scope struct {
	Unrestricted      bool   // true = nessun filtro (es. ADMIN)
	IDs               []uint // ID ammessi (vuoto = nessun risultato, se non Unrestricted)
	FilterByForeignID bool   // true = filtrare per foreign_id (caso AGENCY viste da AGENT)
}

type ScopePolicy interface {
	AgentScope(actor AuthContext) (Scope, error)
	AgencyScope(actor AuthContext) (Scope, error)
}

type scopePolicy struct {
	scopeRepo     repositories.ScopeRepo
	agentNodeRepo repositories.AgentNodeRepo
}

func NewScopePolicy(scopeRepo repositories.ScopeRepo, agentNodeRepo repositories.AgentNodeRepo) ScopePolicy {
	return &scopePolicy{scopeRepo: scopeRepo, agentNodeRepo: agentNodeRepo}
}

func (p *scopePolicy) AgentScope(actor AuthContext) (Scope, error) {
	switch enums.Role(actor.Role) {
	case enums.RoleAdmin:
		return Scope{Unrestricted: true}, nil

	case enums.RoleOperator:
		ids, err := p.scopeRepo.AssignedAgentIDs(actor.UserID)
		return Scope{IDs: ids}, err

	case enums.RoleAgent:
		ids, err := p.agentNodeRepo.GetDescendantAgentIDs(actor.UserID)
		return Scope{IDs: ids}, err

	default:
		return Scope{}, errors.New("ruolo non autorizzato")
	}
}

func (p *scopePolicy) AgencyScope(actor AuthContext) (Scope, error) {
	switch enums.Role(actor.Role) {
	case enums.RoleAdmin:
		return Scope{Unrestricted: true}, nil

	case enums.RoleOperator:
		// operator_agency contiene direttamente gli ID delle AGENCY assegnate
		ids, err := p.scopeRepo.AssignedAgencyIDs(actor.UserID)
		return Scope{IDs: ids}, err

	case enums.RoleAgent:
		// l'agente vede le AGENCY il cui foreign_id appartiene al suo sottoalbero
		ids, err := p.agentNodeRepo.GetDescendantAgentIDs(actor.UserID)
		return Scope{IDs: ids, FilterByForeignID: true}, err

	default:
		return Scope{}, errors.New("ruolo non autorizzato")
	}
}

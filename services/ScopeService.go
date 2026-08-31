package services

import (
	"errors"
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/repositories"
	"fmt"
)

type ScopeService interface {
	AssignToOperator(request pivot.ArraysToOpRequest) (*pivot.ArraysToOpResponse, error)
}

type scopeService struct {
	repo repositories.ScopeRepo
}

func NewScopeService(repo repositories.ScopeRepo) ScopeService {
	return &scopeService{repo: repo}
}

func (s *scopeService) AssignToOperator(request pivot.ArraysToOpRequest) (*pivot.ArraysToOpResponse, error) {
	var agentIds []uint
	if request.AgentIds != nil {
		agentIds = *request.AgentIds
	}

	var agencyIds []uint
	if request.AgencyIds != nil {
		agencyIds = *request.AgencyIds
	}

	// Verifica che gli agent_id forniti formino un sottoalbero "contiguo":
	// non deve mancare nessun nodo compreso tra il primo e l'ultimo (per lft/rgt),
	// altrimenti l'assegnazione spezzerebbe la struttura nested set.
	if err := s.validateContiguousAgentSubtree(agentIds); err != nil {
		return nil, err
	}

	response, err := s.repo.AssignToOperator(request.OperatorId, agentIds, agencyIds)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// validateContiguousAgentSubtree controlla che, presi gli agent_id in input,
// non manchi alcun nodo intermedio nell'intervallo [minLft, maxRgt] che essi definiscono.
func (s *scopeService) validateContiguousAgentSubtree(agentIds []uint) error {
	if len(agentIds) <= 1 {
		return nil
	}

	nodes, err := s.repo.GetNodesByAgentIDs(agentIds)
	if err != nil {
		return err
	}

	if len(nodes) != len(agentIds) {
		return errors.New("uno o più agent_id forniti non corrispondono a nodi esistenti nell'albero")
	}

	minLft, maxRgt := nodes[0].Lft, nodes[0].Rgt
	for _, n := range nodes[1:] {
		if n.Lft < minLft {
			minLft = n.Lft
		}
		if n.Rgt > maxRgt {
			maxRgt = n.Rgt
		}
	}

	idsInRange, err := s.repo.GetAgentIDsInLftRgtRange(minLft, maxRgt)
	if err != nil {
		return err
	}

	provided := make(map[uint]bool, len(agentIds))
	for _, id := range agentIds {
		provided[id] = true
	}

	for _, id := range idsInRange {
		if !provided[id] {
			return fmt.Errorf("selezione non contigua: manca l'agente con id %d, necessario per non spezzare la struttura ad albero", id)
		}
	}

	return nil
}

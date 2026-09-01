package services

import (
	"errors"
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"
	"fmt"
	"reflect"
)

type ScopeService interface {
	AssignToOperator(request pivot.ArraysToOpRequest, actor policies.AuthContext) (*pivot.ArraysToOpResponse, error)
}

type scopeService struct {
	repo       repositories.ScopeRepo
	userRepo   repositories.UserRepo
	logService ActivityLogService
}

func NewScopeService(
	repo repositories.ScopeRepo,
	userRepo repositories.UserRepo,
	logService ActivityLogService,
) ScopeService {
	return &scopeService{
		repo:       repo,
		userRepo:   userRepo,
		logService: logService,
	}
}

func (s *scopeService) AssignToOperator(request pivot.ArraysToOpRequest, actor policies.AuthContext) (*pivot.ArraysToOpResponse, error) {
	var agentIds []uint
	if request.AgentIds != nil {
		agentIds = *request.AgentIds
	}

	var agencyIds []uint
	if request.AgencyIds != nil {
		agencyIds = *request.AgencyIds
	}

	// Verifica che gli agent_id forniti formino un sottoalbero "contiguo"
	if err := s.validateContiguousAgentSubtree(agentIds); err != nil {
		return nil, err
	}

	// Verifica che ogni agenzia fornita appartenga alla catena di agenti selezionata
	if err := s.validateAgenciesBelongToAgents(agencyIds, agentIds); err != nil {
		return nil, err
	}

	response, err := s.repo.AssignToOperator(request.OperatorId, agentIds, agencyIds)
	if err != nil {
		return nil, err
	}

	// Controllo di sicurezza se la response è nil
	if response == nil {
		return nil, errors.New("risposta nulla ritornata dal repository")
	}

	// Helper inline per calcolare in sicurezza la lunghezza di una slice gestita come puntatore
	safeLen := func(slicePtr *[]uint) int {
		if slicePtr == nil {
			return 0
		}
		return len(*slicePtr)
	}

	// Activity Log - Aggiornamento Dati
	err = s.logService.NewLog(&models.ActivityLog{
		ActorID:     &actor.UserID,
		ActorRole:   enums.Role(actor.Role),
		Action:      enums.Assignment,
		TargetType:  reflect.TypeOf(&models.User{}).Elem().Name(), // Restituisce "User" in sicurezza
		TargetID:    &response.OperatorId,
		Description: fmt.Sprintf("Rete assegnata all'operatore %d (Agenti: %d, Agenzie: %d)", request.OperatorId, safeLen(response.AgentIds), safeLen(response.AgencyIds)),
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// validateAgenciesBelongToAgents controlla che ogni agenzia in agencyIds abbia
// ForeignId (l'agente a cui è agganciata) presente tra gli agentIds selezionati.
// Se agencyIds non è vuoto ma agentIds lo è, ogni agenzia fallirà il controllo
// (nessuna catena selezionata a cui appartenere).
func (s *scopeService) validateAgenciesBelongToAgents(agencyIds []uint, agentIds []uint) error {
	if len(agencyIds) == 0 {
		return nil
	}

	agencies, err := s.userRepo.GetAllByRoleAndIDs(enums.RoleAgency.String(), agencyIds)
	if err != nil {
		return err
	}

	if len(agencies) != len(agencyIds) {
		return errors.New("uno o più agency_id forniti non corrispondono ad agenzie esistenti")
	}

	allowedAgents := make(map[uint]bool, len(agentIds))
	for _, id := range agentIds {
		allowedAgents[id] = true
	}

	for _, agency := range agencies {
		if agency.ForeignID == nil || !allowedAgents[*agency.ForeignID] {
			return fmt.Errorf(
				"l'agenzia con id %d non appartiene alla catena di agenti selezionata",
				agency.ID,
			)
		}
	}

	return nil
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

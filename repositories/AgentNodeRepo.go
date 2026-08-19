package repositories

import (
	"errors"
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AgentNodeRepo interface {
	WithTx(tx *gorm.DB) AgentNodeRepo

	GetParentIdByAgentId(agentId uint) (*uint, error)
	Create(nodeModel *models.AgentNode) error
	GetAncestors(id *uint) ([]models.AgentNode, error)
	GetDescendants(id *uint) ([]models.AgentNode, error)
	GetTrees() ([]*models.AgentNode, error)
	GetFilteredTreeByAgent(userID uint) ([]*models.AgentNode, error)
	GetFilteredTreeByOperator(userID uint) ([]*models.AgentNode, error)

	GetNodeByAgentID(agentID uint) (*models.AgentNode, error)
	IsDescendantOrSelf(ancestorAgentID uint, nodeID uint) (bool, error)
	IsAgentDescendantOrSelf(ancestorAgentID, candidateAgentID uint) (bool, error)
	GetDescendantAgentIDs(ancestorAgentID uint) ([]uint, error)
}

type agentNodeRepo struct {
	db *gorm.DB
}

func NewAgentNodeRepository(db *gorm.DB) AgentNodeRepo {
	return &agentNodeRepo{db: db}
}

func (r *agentNodeRepo) WithTx(tx *gorm.DB) AgentNodeRepo {
	return &agentNodeRepo{db: tx}
}

func (r *agentNodeRepo) GetParentIdByAgentId(agentId uint) (*uint, error) {
	var parentId *uint

	err := r.db.Model(&models.AgentNode{}).
		Where("agent_id = ?", agentId).
		Pluck("parent_id", &parentId).Error

	if err != nil {
		return nil, err
	}
	return parentId, nil
}

// Creazione dell'agente e assegnazione al nodo di appartenenza
// Aggiornamento dei campi lft e rgt di tutto l'albero secondo le regole della struttura nested
func (r *agentNodeRepo) Create(nodeModel *models.AgentNode) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if nodeModel.Agent != nil {
			if err := tx.Create(nodeModel.Agent).Error; err != nil {
				return err
			}
			nodeModel.AgentID = nodeModel.Agent.ID
		}

		if nodeModel.ParentID == nil {
			var lastRoot models.AgentNode

			err := tx.
				Where("parent_id IS NULL").
				Order("rgt DESC").
				First(&lastRoot).Error

			switch {
			case err == nil:
				// Esiste già almeno una root: agganciati alla sua rgt
				nodeModel.Lft = lastRoot.Rgt + 1
				nodeModel.Rgt = nodeModel.Lft + 1

			case errors.Is(err, gorm.ErrRecordNotFound):
				// Prima root in assoluto
				nodeModel.Lft = 1
				nodeModel.Rgt = 2

			default:
				return err
			}
			return tx.Omit("Agent", "Parent").Create(nodeModel).Error
		}

		var parent models.AgentNode
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&parent, *nodeModel.ParentID).Error; err != nil {
			return err
		}

		parentRgt := parent.Rgt

		if err := tx.Model(&models.AgentNode{}).
			Where("rgt >= ? OR lft > ?", parentRgt, parentRgt).
			Updates(map[string]interface{}{
				"rgt": gorm.Expr("CASE WHEN rgt >= ? THEN rgt + 2 ELSE rgt END", parentRgt),
				"lft": gorm.Expr("CASE WHEN lft > ? THEN lft + 2 ELSE lft END", parentRgt),
			}).Error; err != nil {
			return err
		}

		nodeModel.Lft = parentRgt
		nodeModel.Rgt = parentRgt + 1

		return tx.Omit("Agent", "Parent").Create(nodeModel).Error
	})
}

func (r *agentNodeRepo) GetAncestors(id *uint) ([]models.AgentNode, error) {
	if id == nil {
		return nil, fmt.Errorf("id non può essere nil")
	}

	var child models.AgentNode
	if err := r.db.First(&child, *id).Error; err != nil {
		return nil, err
	}

	var ancestors []models.AgentNode
	err := r.db.Where("lft < ? AND rgt > ?", child.Lft, child.Rgt).
		Order("lft ASC").
		Find(&ancestors).Error

	return ancestors, err
}

func (r *agentNodeRepo) GetDescendants(id *uint) ([]models.AgentNode, error) {
	if id == nil {
		return nil, fmt.Errorf("id non può essere nil")
	}

	var parent models.AgentNode
	if err := r.db.First(&parent, *id).Error; err != nil {
		return nil, err
	}

	var descendants []models.AgentNode
	err := r.db.Where("lft > ? AND rgt < ?", parent.Lft, parent.Rgt).
		Order("lft ASC").
		Find(&descendants).Error

	return descendants, err
}

func (r *agentNodeRepo) GetTrees() ([]*models.AgentNode, error) {
	var nodes []*models.AgentNode

	// 1. Recupera tutti i nodi ordinati per 'lft' con un'unica query.
	err := r.db.Preload("Agent").Order("lft ASC").Find(&nodes).Error
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return []*models.AgentNode{}, nil
	}

	// 1bis. Raccoglie gli AgentID unici per recuperare le agenzie in un'unica query
	agentIDs := make([]uint, 0, len(nodes))
	seenAgentID := make(map[uint]bool, len(nodes))
	for _, node := range nodes {
		if node.AgentID != 0 && !seenAgentID[node.AgentID] {
			seenAgentID[node.AgentID] = true
			agentIDs = append(agentIDs, node.AgentID)
		}
	}

	agenciesByAgent := make(map[uint][]*models.User)
	if len(agentIDs) > 0 {
		var agencies []*models.User
		if err := r.db.
			Where("role = ? AND foreign_id IN ?", enums.RoleAgency, agentIDs).
			Find(&agencies).Error; err != nil {
			return nil, err
		}

		for _, agency := range agencies {
			if agency.ForeignId == nil {
				continue
			}
			key := *agency.ForeignId
			agenciesByAgent[key] = append(agenciesByAgent[key], agency)
		}
	}

	// 2. Un solo passaggio: grazie all'ordinamento 'lft ASC' il genitore di un
	// nodo si trova SEMPRE prima del nodo stesso nella slice (proprietà del
	// nested set), quindi quando processiamo "node" il suo Parent è già in
	// nodeMap e possiamo collegarlo subito, senza un secondo ciclo separato.
	var roots []*models.AgentNode
	nodeMap := make(map[uint]*models.AgentNode, len(nodes))

	for _, node := range nodes {
		node.Children = make([]*models.AgentNode, 0)

		if agencies, ok := agenciesByAgent[node.AgentID]; ok {
			node.Agencies = agencies
		} else {
			node.Agencies = make([]*models.User, 0)
		}

		nodeMap[node.ID] = node
		if node.ParentID == nil {
			roots = append(roots, node)
			continue
		}

		if parent, exists := nodeMap[*node.ParentID]; exists {
			parent.Children = append(parent.Children, node)
		}
	}

	return roots, nil
}

// -----------------------------------------------------------------------------
// LOGICA AGENT: Ritorna l'agente stesso + tutti i figli e relative agenzie
// -----------------------------------------------------------------------------
func (r *agentNodeRepo) GetFilteredTreeByAgent(userID uint) ([]*models.AgentNode, error) {
	// 1. Trova il nodo radice dell'agente loggato
	var targetNode models.AgentNode
	if err := r.db.Where("agent_id = ?", userID).First(&targetNode).Error; err != nil {
		return nil, fmt.Errorf("nodo agente non trovato per userID %d: %w", userID, err)
	}

	// 2. Sfrutta il Nested Set: tutti i discendenti hanno lft/rgt compresi nel nodo padre
	var nodes []*models.AgentNode
	err := r.db.Preload("Agent").
		Where("lft >= ? AND rgt <= ?", targetNode.Lft, targetNode.Rgt).
		Order("lft ASC").
		Find(&nodes).Error

	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return []*models.AgentNode{}, nil
	}

	// 3. Costruisce l'albero e associa le agenzie
	return r.buildTreeWithAgencies(nodes, nil)
}

// -----------------------------------------------------------------------------
// LOGICA OPERATOR: Ritorna solo agenti ed agenzie presenti nelle tabelle pivot
// -----------------------------------------------------------------------------
func (r *agentNodeRepo) GetFilteredTreeByOperator(operatorID uint) ([]*models.AgentNode, error) {
	// 1. Recupera gli agent_id associati all'operatore dalla tabella pivot agent_operator
	var allowedAgentIDs []uint
	if err := r.db.Table("agent_operator").
		Where("operator_id = ?", operatorID).
		Pluck("agent_id", &allowedAgentIDs).Error; err != nil {
		return nil, err
	}

	if len(allowedAgentIDs) == 0 {
		return []*models.AgentNode{}, nil
	}

	// 2. Recupera le agenzie (ID utente o ForeignID) abilitate dalla tabella pivot agency_operator
	var allowedAgencyIDs []uint
	if err := r.db.Table("agency_operator").
		Where("operator_id = ?", operatorID).
		Pluck("agency_id", &allowedAgencyIDs).Error; err != nil {
		return nil, err
	}

	// 3. Recupera solo i nodi degli agenti associati all'operatore
	var nodes []*models.AgentNode
	err := r.db.Preload("Agent").
		Where("agent_id IN ?", allowedAgentIDs).
		Order("lft ASC").
		Find(&nodes).Error

	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return []*models.AgentNode{}, nil
	}

	// 4. Costruisce l'albero applicando la restrizione per agenzie consentite
	return r.buildTreeWithAgencies(nodes, allowedAgencyIDs)
}

// -----------------------------------------------------------------------------
// HELPER: Funzione di utilità per caricare le agenzie e assemblare l'albero
// -----------------------------------------------------------------------------
func (r *agentNodeRepo) buildTreeWithAgencies(nodes []*models.AgentNode, allowedAgencyIDs []uint) ([]*models.AgentNode, error) {
	// 1. Estrae i vari agentID univoci dai nodi recuperati
	agentIDs := make([]uint, 0, len(nodes))
	seenAgentID := make(map[uint]bool, len(nodes))
	for _, node := range nodes {
		if node.AgentID != 0 && !seenAgentID[node.AgentID] {
			seenAgentID[node.AgentID] = true
			agentIDs = append(agentIDs, node.AgentID)
		}
	}

	agenciesByAgent := make(map[uint][]*models.User)
	if len(agentIDs) > 0 {
		var agencies []*models.User
		query := r.db.Where("role = ? AND foreign_id IN ?", enums.RoleAgency, agentIDs)

		// Se stiamo filtrando per operatore (allowedAgencyIDs != nil), applica la restrizione sulle agenzie
		if allowedAgencyIDs != nil {
			if len(allowedAgencyIDs) == 0 {
				query = query.Where("1 = 0") // Nessuna agenzia visibile
			} else {
				query = query.Where("id IN ?", allowedAgencyIDs)
			}
		}

		if err := query.Find(&agencies).Error; err != nil {
			return nil, err
		}

		for _, agency := range agencies {
			if agency.ForeignId == nil {
				continue
			}
			key := *agency.ForeignId
			agenciesByAgent[key] = append(agenciesByAgent[key], agency)
		}
	}

	// 2. Ricostruzione dell'albero in un singolo ciclo
	var roots []*models.AgentNode
	nodeMap := make(map[uint]*models.AgentNode, len(nodes))

	for _, node := range nodes {
		node.Children = make([]*models.AgentNode, 0)

		if agencies, ok := agenciesByAgent[node.AgentID]; ok {
			node.Agencies = agencies
		} else {
			node.Agencies = make([]*models.User, 0)
		}

		nodeMap[node.ID] = node

		// Se il padre non esiste in nodeMap (perché è la radice del sottoalbero o è stato filtrato via),
		// questo nodo diventa un nodo principale ("root") nel nostro risultato.
		if node.ParentID == nil {
			roots = append(roots, node)
		} else if parent, exists := nodeMap[*node.ParentID]; exists {
			parent.Children = append(parent.Children, node)
		} else {
			roots = append(roots, node)
		}
	}

	return roots, nil
}

// GetNodeByAgentID recupera il nodo dell'albero corrispondente a un dato AgentID (User.ID)
func (r *agentNodeRepo) GetNodeByAgentID(agentID uint) (*models.AgentNode, error) {
	var node models.AgentNode
	if err := r.db.Where("agent_id = ?", agentID).First(&node).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

// IsDescendantOrSelf verifica se il nodo con ID = nodeID è lo stesso nodo dell'agente
// ancestorAgentID, oppure un suo discendente nel nested set.
// Usato per validare il ParentId quando un AGENT crea un nuovo AgentNode.
func (r *agentNodeRepo) IsDescendantOrSelf(ancestorAgentID uint, nodeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.AgentNode{}).
		Where(`id = ? AND (
			id = (SELECT id FROM agent_nodes WHERE agent_id = ?)
			OR (lft > (SELECT lft FROM agent_nodes WHERE agent_id = ?)
			    AND rgt < (SELECT rgt FROM agent_nodes WHERE agent_id = ?))
		)`, nodeID, ancestorAgentID, ancestorAgentID, ancestorAgentID).
		Count(&count).Error
	return count > 0, err
}

// IsAgentDescendantOrSelf come sopra, ma prendendo in input due AgentID (User.ID)
// invece di un AgentNode.ID. Usato per validare il ForeignId quando un AGENT
// crea una nuova AGENCY collegata a un agente del proprio sottoalbero.
func (r *agentNodeRepo) IsAgentDescendantOrSelf(ancestorAgentID, candidateAgentID uint) (bool, error) {
	if ancestorAgentID == candidateAgentID {
		return true, nil
	}

	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) FROM agent_nodes candidate
		JOIN agent_nodes ancestor ON ancestor.agent_id = ?
		WHERE candidate.agent_id = ?
		  AND candidate.lft > ancestor.lft
		  AND candidate.rgt < ancestor.rgt
	`, ancestorAgentID, candidateAgentID).Scan(&count).Error

	return count > 0, err
}

// GetDescendantAgentIDs ritorna gli AgentID (User.ID) di tutti i nodi del sottoalbero
// dell'agente indicato, incluso se stesso. Utile per filtrare liste (agenti/agenzie visibili).
func (r *agentNodeRepo) GetDescendantAgentIDs(ancestorAgentID uint) ([]uint, error) {
	ancestorNode, err := r.GetNodeByAgentID(ancestorAgentID)
	if err != nil {
		return nil, err
	}

	var ids []uint
	err = r.db.Model(&models.AgentNode{}).
		Where("(lft > ? AND rgt < ?) OR id = ?", ancestorNode.Lft, ancestorNode.Rgt, ancestorNode.ID).
		Pluck("agent_id", &ids).Error
	return ids, err
}

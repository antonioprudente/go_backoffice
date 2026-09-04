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

	Create(nodeModel *models.AgentNode) error
	GetTrees() ([]*models.AgentNode, error)
	GetFilteredTreeByAgent(userID uint) ([]*models.AgentNode, error)
	GetFilteredTreeByOperator(userID uint) ([]*models.AgentNode, error)
	GetNodeByAgentID(agentID uint) (*models.AgentNode, error)
	DeleteAgentNodeAndAgentByAgentID(agentID uint) error
	RestoreAgentSubtree(rootAgentID uint) error
	MoveNode(agentID uint, ForeignID *uint) error
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

// GetTrees recupera l'intero albero (nessuna restrizione su agenti/agenzie).
// La costruzione dell'albero e l'associazione delle agenzie sono delegate a
// buildTreeWithAgencies (allowedAgencyIDs = nil -> nessun filtro sulle agenzie),
// evitando di duplicare qui la stessa logica già usata dalle viste filtrate.
func (r *agentNodeRepo) GetTrees() ([]*models.AgentNode, error) {
	var nodes []*models.AgentNode
	if err := r.db.Preload("Agent").Order("lft ASC").Find(&nodes).Error; err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return []*models.AgentNode{}, nil
	}

	return r.buildTreeWithAgencies(nodes, nil)
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
			if agency.ForeignID == nil {
				continue
			}
			key := *agency.ForeignID
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

func (r *agentNodeRepo) DeleteAgentNodeAndAgentByAgentID(agentID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var node models.AgentNode
		// Lock per evitare che un altro Create/Delete concorrente modifichi
		// lft/rgt mentre stiamo calcolando il subtree da rimuovere
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("agent_id = ?", agentID).
			First(&node).Error; err != nil {
			return err
		}

		width := node.Rgt - node.Lft + 1

		// Recupera l'intero subtree ordinato per lft DESC: grazie alla proprietà
		// del nested set (child.lft > parent.lft sempre), questo ordine garantisce
		// che ogni figlio venga elaborato/cancellato prima del proprio genitore
		var subtreeNodes []models.AgentNode
		if err := tx.
			Where("lft >= ? AND rgt <= ?", node.Lft, node.Rgt).
			Order("lft DESC").
			Find(&subtreeNodes).Error; err != nil {
			return err
		}

		subtreeAgentIDs := make([]uint, 0, len(subtreeNodes))
		subtreeNodeIDs := make([]uint, 0, len(subtreeNodes))
		for _, n := range subtreeNodes {
			subtreeAgentIDs = append(subtreeAgentIDs, n.AgentID)
			subtreeNodeIDs = append(subtreeNodeIDs, n.ID)
		}

		if len(subtreeAgentIDs) > 0 {
			var agencyIDs []uint
			if err := tx.Model(&models.User{}).
				Where("role = ? AND foreign_id IN ?", enums.RoleAgency, subtreeAgentIDs).
				Pluck("id", &agencyIDs).Error; err != nil {
				return err
			}

			if len(agencyIDs) > 0 {
				// Utenti finali agganciati alle agenzie del sottoalbero -> soft delete (deleted_at)
				if err := tx.Where("role = ? AND foreign_id IN ?", enums.RoleUser, agencyIDs).
					Delete(&models.User{}).Error; err != nil {
					return err
				}
				// Agenzie del sottoalbero -> soft delete (deleted_at)
				if err := tx.Where("id IN ?", agencyIDs).
					Delete(&models.User{}).Error; err != nil {
					return err
				}
			}

			// Agenti del sottoalbero (incluso quello radice) -> soft delete (deleted_at)
			if err := tx.Where("role = ? AND id IN ?", enums.RoleAgent, subtreeAgentIDs).
				Delete(&models.User{}).Error; err != nil {
				return err
			}
		}

		// Hard delete dei nodi in ordine lft DESC (figli prima dei genitori),
		// per rispettare il vincolo FK auto-referenziale fk_agent_nodes_parent
		for _, nodeID := range subtreeNodeIDs {
			if err := tx.Where("id = ?", nodeID).
				Delete(&models.AgentNode{}).Error; err != nil {
				return err
			}
		}

		// Richiude il "buco" lasciato dal sottoalbero eliminato nel nested set
		if err := tx.Model(&models.AgentNode{}).
			Where("rgt > ?", node.Rgt).
			Update("rgt", gorm.Expr("rgt - ?", width)).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.AgentNode{}).
			Where("lft > ?", node.Rgt).
			Update("lft", gorm.Expr("lft - ?", width)).Error; err != nil {
			return err
		}

		return nil
	})
}

// RestoreAgentSubtree ripristina un agente eliminato (e tutto il suo sottoalbero)
// ricostruendo gli AgentNode a partire dai ForeignID rimasti sugli User soft-deleted.
// NB: le posizioni esatte lft/rgt e l'ordine tra fratelli non sono garantiti identici
// all'originale, ma la struttura genitore-figlio viene ripristinata fedelmente.
func (r *agentNodeRepo) RestoreAgentSubtree(rootAgentID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		nodeRepo := r.WithTx(tx)

		var rootUser models.User
		if err := tx.Unscoped().
			Where("id = ? AND role = ? AND deleted_at IS NOT NULL", rootAgentID, enums.RoleAgent).
			First(&rootUser).Error; err != nil {
			return fmt.Errorf("agente eliminato non trovato: %w", err)
		}

		// BFS: raccoglie l'intero sottoalbero di agenti eliminati seguendo
		// la catena dei ForeignID (parent -> children). L'ordine di visita
		// garantisce che ogni genitore venga elaborato prima dei suoi figli.
		subtreeAgents := []models.User{rootUser}
		queue := []uint{rootUser.ID}
		for len(queue) > 0 {
			currentID := queue[0]
			queue = queue[1:]

			var children []models.User
			if err := tx.Unscoped().
				Where("role = ? AND foreign_id = ? AND deleted_at IS NOT NULL", enums.RoleAgent, currentID).
				Find(&children).Error; err != nil {
				return err
			}
			for _, child := range children {
				subtreeAgents = append(subtreeAgents, child)
				queue = append(queue, child.ID)
			}
		}

		agentIDs := make([]uint, 0, len(subtreeAgents))
		for _, a := range subtreeAgents {
			agentIDs = append(agentIDs, a.ID)
		}

		// Agenzie eliminate collegate agli agenti del sottoalbero
		var agencyIDs []uint
		if err := tx.Unscoped().Model(&models.User{}).
			Where("role = ? AND foreign_id IN ? AND deleted_at IS NOT NULL", enums.RoleAgency, agentIDs).
			Pluck("id", &agencyIDs).Error; err != nil {
			return err
		}

		// Ripristina gli agenti
		if err := tx.Unscoped().Model(&models.User{}).
			Where("id IN ?", agentIDs).
			Update("deleted_at", nil).Error; err != nil {
			return err
		}

		// Ripristina agenzie e relativi utenti finali
		if len(agencyIDs) > 0 {
			if err := tx.Unscoped().Model(&models.User{}).
				Where("id IN ?", agencyIDs).
				Update("deleted_at", nil).Error; err != nil {
				return err
			}
			if err := tx.Unscoped().Model(&models.User{}).
				Where("role = ? AND foreign_id IN ?", enums.RoleUser, agencyIDs).
				Update("deleted_at", nil).Error; err != nil {
				return err
			}
		}

		// Ricrea gli AgentNode in ordine top-down, riusando l'algoritmo
		// di inserimento nested-set già usato in fase di creazione normale
		for _, agent := range subtreeAgents {
			// Guard: se il nodo esiste già (es. undo chiamato due volte), salta
			if _, err := nodeRepo.GetNodeByAgentID(agent.ID); err == nil {
				continue
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			var parentNodeID *uint
			if agent.ForeignID != nil {
				parentNode, err := nodeRepo.GetNodeByAgentID(*agent.ForeignID)
				if err != nil {
					return fmt.Errorf("nodo padre non trovato per l'agente %d: %w", agent.ID, err)
				}
				parentNodeID = &parentNode.ID
			}

			newNode := &models.AgentNode{
				ParentID: parentNodeID,
				AgentID:  agent.ID,
			}
			if err := nodeRepo.Create(newNode); err != nil {
				return fmt.Errorf("impossibile ricreare il nodo per l'agente %d: %w", agent.ID, err)
			}
		}

		return nil
	})
}

// MoveNode sposta il nodo identificato da agentID (e tutto il suo sottoalbero).
// Se ForeignID è nil, il nodo diventa una nuova radice (root); altrimenti viene
// spostato come ultimo figlio del nodo il cui agent_id è *ForeignID.
// Aggiorna lft/rgt di tutta la tabella secondo le regole del nested set
// e il ParentID del nodo spostato.
func (r *agentNodeRepo) MoveNode(agentID uint, ForeignID *uint) error {
	if ForeignID != nil && agentID == *ForeignID {
		return fmt.Errorf("impossibile spostare il nodo %d sotto se stesso", agentID)
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		var node models.AgentNode
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("agent_id = ?", agentID).
			First(&node).Error; err != nil {
			return err
		}

		var newParentNode *models.AgentNode
		if ForeignID != nil {
			var np models.AgentNode
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("agent_id = ?", *ForeignID).
				First(&np).Error; err != nil {
				return err
			}
			if np.Lft >= node.Lft && np.Rgt <= node.Rgt {
				return fmt.Errorf("impossibile spostare il nodo %d sotto se stesso o un suo discendente (foreign_id %d)", agentID, *ForeignID)
			}
			newParentNode = &np
		}

		// Update("colonna", valore) scrive il valore così com'è, incluso nil -> NULL,
		// a differenza di Updates(struct{}) che salterebbe i puntatori nil.
		if err := tx.Model(&models.User{}).
			Where("id = ?", agentID).
			Update("foreign_id", ForeignID).Error; err != nil {
			return err
		}

		// Cattura l'intero sottoalbero (radice + discendenti), ordinato per lft
		// così l'ordine relativo tra fratelli viene preservato al reinserimento
		var subtree []models.AgentNode
		if err := tx.Where("lft >= ? AND rgt <= ?", node.Lft, node.Rgt).
			Order("lft ASC").
			Find(&subtree).Error; err != nil {
			return err
		}

		byID := make(map[uint]*models.AgentNode, len(subtree))
		for i := range subtree {
			byID[subtree[i].ID] = &subtree[i]
		}
		// Per ogni discendente (radice esclusa), l'AgentID del suo genitore
		// ALL'INTERNO del sottoalbero
		parentAgentIDOf := make(map[uint]uint, len(subtree))
		for i := range subtree {
			n := &subtree[i]
			if n.ID == node.ID {
				continue
			}
			parentAgentIDOf[n.AgentID] = byID[*n.ParentID].AgentID
		}

		// 1) elimina l'intero sottoalbero e richiude il buco
		if err := r.deleteSubtree(tx, &node); err != nil {
			return err
		}

		// 2) reinserisce la radice del sottoalbero sotto la nuova destinazione
		var newParentNodeID *uint
		if newParentNode != nil {
			newParentNodeID = &newParentNode.ID
		}
		newRoot, err := r.insertNodeUnderParent(tx, node.AgentID, newParentNodeID)
		if err != nil {
			return err
		}

		// 3) reinserisce ricorsivamente ogni discendente, nello stesso ordine
		//    relativo originale, sotto il proprio genitore già reinserito
		newNodeIDByAgentID := map[uint]uint{node.AgentID: newRoot.ID}
		for _, n := range subtree {
			if n.ID == node.ID {
				continue
			}
			parentNewID, ok := newNodeIDByAgentID[parentAgentIDOf[n.AgentID]]
			if !ok {
				return fmt.Errorf("errore interno: genitore di agent_id %d non ancora reinserito", n.AgentID)
			}
			inserted, err := r.insertNodeUnderParent(tx, n.AgentID, &parentNewID)
			if err != nil {
				return err
			}
			newNodeIDByAgentID[n.AgentID] = inserted.ID
		}

		return nil
	})
}

// insertNodeUnderParent inserisce (o reinserisce) il nodo con il dato agentID
// come ultimo figlio del nodo il cui ID è parentNodeID.
// Se parentNodeID è nil, il nodo diventa una nuova radice: il ParentID DEVE
// essere scritto come NULL, non semplicemente "non toccato". Con Updates(struct)
// GORM ignora i campi a valore zero (incluso un puntatore nil), quindi qui
// usiamo Create (che scrive sempre tutte le colonne) e, in ogni punto in cui
// in futuro si volesse fare un Update, va sempre usato Select("ParentID")
// o una map esplicita per garantire la scrittura del NULL.
func (r *agentNodeRepo) insertNodeUnderParent(tx *gorm.DB, agentID uint, parentNodeID *uint) (*models.AgentNode, error) {
	newNode := &models.AgentNode{
		AgentID:  agentID,
		ParentID: parentNodeID, // può essere nil: va bene, Create scrive sempre la colonna
	}

	if parentNodeID == nil {
		var lastRoot models.AgentNode
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("parent_id IS NULL").
			Order("rgt DESC").
			First(&lastRoot).Error

		switch {
		case err == nil:
			newNode.Lft = lastRoot.Rgt + 1
			newNode.Rgt = newNode.Lft + 1
		case errors.Is(err, gorm.ErrRecordNotFound):
			newNode.Lft = 1
			newNode.Rgt = 2
		default:
			return nil, err
		}

		if err := tx.Create(newNode).Error; err != nil {
			return nil, err
		}
		return newNode, nil
	}

	var parent models.AgentNode
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&parent, *parentNodeID).Error; err != nil {
		return nil, err
	}

	parentRgt := parent.Rgt

	if err := tx.Model(&models.AgentNode{}).
		Where("rgt >= ? OR lft > ?", parentRgt, parentRgt).
		Updates(map[string]interface{}{
			"rgt": gorm.Expr("CASE WHEN rgt >= ? THEN rgt + 2 ELSE rgt END", parentRgt),
			"lft": gorm.Expr("CASE WHEN lft > ? THEN lft + 2 ELSE lft END", parentRgt),
		}).Error; err != nil {
		return nil, err
	}

	newNode.Lft = parentRgt
	newNode.Rgt = parentRgt + 1

	if err := tx.Create(newNode).Error; err != nil {
		return nil, err
	}
	return newNode, nil
}

// deleteSubtree rimuove il nodo e tutti i suoi discendenti, richiudendo
// il buco lasciato nella struttura nested set per il resto dell'albero.
func (r *agentNodeRepo) deleteSubtree(tx *gorm.DB, node *models.AgentNode) error {
	width := node.Rgt - node.Lft + 1

	// 1. Elimina i nodi in ordine decrescente di LFT (prima le foglie/figli, poi i padri)
	if err := tx.Where("lft >= ? AND rgt <= ?", node.Lft, node.Rgt).
		Order("lft DESC").
		Delete(&models.AgentNode{}).Error; err != nil {
		return err
	}

	// 2. Aggiorna i valori RGT dei nodi rimanenti
	if err := tx.Model(&models.AgentNode{}).
		Where("rgt > ?", node.Rgt).
		Update("rgt", gorm.Expr("rgt - ?", width)).Error; err != nil {
		return err
	}

	// 3. Aggiorna i valori LFT dei nodi rimanenti
	if err := tx.Model(&models.AgentNode{}).
		Where("lft > ?", node.Rgt).
		Update("lft", gorm.Expr("lft - ?", width)).Error; err != nil {
		return err
	}

	return nil
}

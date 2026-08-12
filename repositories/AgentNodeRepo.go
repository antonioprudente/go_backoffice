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
	Create(nodeModel *models.AgentNode) error
	GetAncestors(id *uint) ([]models.AgentNode, error)
	GetDescendants(id *uint) ([]models.AgentNode, error)
	GetTrees() ([]*models.AgentNode, error)
}

type agentNodeRepo struct {
	db *gorm.DB
}

func NewAgentNodeRepository(db *gorm.DB) AgentNodeRepo {
	return &agentNodeRepo{db: db}
}

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
			Where("rgt >= ?", parentRgt).
			Update("rgt", gorm.Expr("rgt + 2")).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.AgentNode{}).
			Where("lft > ?", parentRgt).
			Update("lft", gorm.Expr("lft + 2")).Error; err != nil {
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
	// Precarica l'Agent se presente.
	err := r.db.Preload("Agent").Order("lft ASC").Find(&nodes).Error
	if err != nil {
		return nil, err
	}

	// Se il database è vuoto, ritorna una lista vuota
	if len(nodes) == 0 {
		return []*models.AgentNode{}, nil
	}

	// 1bis. Raccoglie tutti gli AgentID presenti nei nodi, per recuperare
	// in un'unica query tutte le agenzie collegate (User con role = AGENCY
	// e foreign_id = AgentID di uno degli agenti dell'albero)
	agentIDs := make([]uint, 0, len(nodes))
	seenAgentID := make(map[uint]bool)
	for _, node := range nodes {
		if node.AgentID != 0 && !seenAgentID[node.AgentID] {
			seenAgentID[node.AgentID] = true
			agentIDs = append(agentIDs, node.AgentID)
		}
	}

	// Mappa agentID -> slice di agenzie (User con Role=AGENCY e ForeignId == agentID)
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

	var roots []*models.AgentNode

	// Mappa di supporto per rintracciare i nodi genitori istantaneamente tramite ID
	nodeMap := make(map[uint]*models.AgentNode)

	// 2. Primo passaggio: mappa tutti i nodi con il loro ID e assegna le agenzie
	for _, node := range nodes {
		// Inizializza la fetta dei figli per evitare slice nil nel JSON
		node.Children = make([]*models.AgentNode, 0)
		// Assegna le agenzie (role=AGENCY) collegate a questo agente, se presenti

		if agencies, ok := agenciesByAgent[node.AgentID]; ok {
			node.Agencies = agencies
		} else {
			node.Agencies = make([]*models.User, 0)
		}
		nodeMap[node.ID] = node
	}

	// 3. Secondo passaggio: costruisce l'albero collegando i figli ai genitori
	for _, node := range nodes {
		if node.ParentID == nil {
			// È un nodo radice (può essercene più di uno se usi il multi-tree)
			roots = append(roots, node)
		} else {
			// Trova il genitore nella mappa e appendi il nodo corrente ai suoi figli
			if parent, exists := nodeMap[*node.ParentID]; exists {
				parent.Children = append(parent.Children, node)
			}
			// Nota: l'ordine 'lft ASC' garantisce che i figli vengano
			// appesi nell'ordine gerarchico e cronologico corretto di lettura.
		}
	}
	return roots, nil
}

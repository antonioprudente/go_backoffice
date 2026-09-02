package repositories

import (
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type ScopeRepo interface {
	AssignToOperator(operatorID uint, agentIds []uint, agencyIds []uint) (*pivot.ArraysToOpResponse, error)

	IsAgentAssignedToOperator(operatorID uint, agentID uint) (bool, error)
	IsAgencyAssignedToOperator(operatorID uint, agencyID uint) (bool, error)
	AssignedAgentIDs(operatorID uint) ([]uint, error)
	AssignedAgencyIDs(operatorID uint) ([]uint, error)
	NodeChildrenAndSelfAgentIds(agentID uint) ([]uint, error)
	NodeChildrenAgentIds(agentID uint) ([]uint, error)

	// Nuovi metodi per la validazione di contiguità del sottoalbero
	GetNodesByAgentIDs(agentIds []uint) ([]models.AgentNode, error)
	GetAgentIDsInLftRgtRange(minLft, maxRgt uint) ([]uint, error)
}

type scopeRepo struct {
	db *gorm.DB
}

func NewScopeRepository(db *gorm.DB) ScopeRepo {
	return &scopeRepo{db: db}
}

func (r *scopeRepo) AssignToOperator(operatorID uint, agentIds []uint, agencyIds []uint) (*pivot.ArraysToOpResponse, error) {
	response := &pivot.ArraysToOpResponse{
		OperatorId: operatorID,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var operator models.User
		if err := tx.Where("id = ?", operatorID).First(&operator).Error; err != nil {
			return err
		}

		// --- Elimina le vecchie associazioni AGENTI per questo operatore ---
		if err := tx.Where("operator_id = ?", operatorID).
			Delete(&models.AgentOperator{}).Error; err != nil {
			return err
		}

		// --- Elimina le vecchie associazioni AGENZIE per questo operatore ---
		if err := tx.Where("operator_id = ?", operatorID).
			Delete(&models.AgencyOperator{}).Error; err != nil {
			return err
		}

		// --- Inserisce le nuove associazioni AGENTI ---
		if len(agentIds) > 0 {
			agentOperators := make([]models.AgentOperator, 0, len(agentIds))
			for _, agentId := range agentIds {
				agentOperators = append(agentOperators, models.AgentOperator{
					OperatorID: operatorID,
					AgentID:    agentId,
				})
			}
			if err := tx.Create(&agentOperators).Error; err != nil {
				return err
			}
			response.AgentIds = &agentIds
		}

		// --- Inserisce le nuove associazioni AGENZIE ---
		if len(agencyIds) > 0 {
			agencyOperators := make([]models.AgencyOperator, 0, len(agencyIds))
			for _, agencyId := range agencyIds {
				agencyOperators = append(agencyOperators, models.AgencyOperator{
					OperatorID: operatorID,
					AgencyID:   agencyId,
				})
			}
			if err := tx.Create(&agencyOperators).Error; err != nil {
				return err
			}
			response.AgencyIds = &agencyIds
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *scopeRepo) IsAgentAssignedToOperator(operatorID uint, agentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.AgentOperator{}).
		Where("operator_id = ? AND agent_id = ?", operatorID, agentID).
		Count(&count).Error
	return count > 0, err
}

func (r *scopeRepo) IsAgencyAssignedToOperator(operatorID uint, agencyID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.AgencyOperator{}).
		Where("operator_id = ? AND agency_id = ?", operatorID, agencyID).
		Count(&count).Error
	return count > 0, err
}

func (r *scopeRepo) AssignedAgentIDs(operatorID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&models.AgentOperator{}).
		Where("operator_id = ?", operatorID).
		Pluck("agent_id", &ids).Error
	return ids, err
}

func (r *scopeRepo) AssignedAgencyIDs(operatorID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&models.AgencyOperator{}).
		Where("operator_id = ?", operatorID).
		Pluck("agency_id", &ids).Error
	return ids, err
}

func (r *scopeRepo) NodeChildrenAndSelfAgentIds(agentID uint) ([]uint, error) {
	var node models.AgentNode
	if err := r.db.Where("agent_id = ?", agentID).First(&node).Error; err != nil {
		return nil, err
	}

	var children []uint
	err := r.db.Model(&models.AgentNode{}).
		Where("lft >= ? AND rgt <= ?", node.Lft, node.Rgt).
		Pluck("agent_id", &children).Error
	if err != nil {
		return nil, err
	}

	return children, nil
}

func (r *scopeRepo) NodeChildrenAgentIds(agentID uint) ([]uint, error) {
	var node models.AgentNode
	if err := r.db.Where("agent_id = ?", agentID).First(&node).Error; err != nil {
		return nil, err
	}

	var children []uint
	err := r.db.Model(&models.AgentNode{}).
		Where("lft > ? AND rgt < ?", node.Lft, node.Rgt).
		Pluck("agent_id", &children).Error
	if err != nil {
		return nil, err
	}

	return children, nil
}

// GetNodesByAgentIDs recupera i nodi (con lft/rgt) corrispondenti agli agent_id forniti.
// Se un id non corrisponde a nessun nodo, semplicemente non comparirà nel risultato
// (il chiamante deve controllare che len(risultato) == len(agentIds)).
func (r *scopeRepo) GetNodesByAgentIDs(agentIds []uint) ([]models.AgentNode, error) {
	var nodes []models.AgentNode
	if len(agentIds) == 0 {
		return nodes, nil
	}
	err := r.db.Where("agent_id IN ?", agentIds).Find(&nodes).Error
	return nodes, err
}

// GetAgentIDsInLftRgtRange ritorna gli agent_id di tutti i nodi il cui lft/rgt
// è interamente contenuto nell'intervallo [minLft, maxRgt]. Usato per verificare
// che un sottoinsieme di nodi non lasci "buchi" nella struttura ad albero.
func (r *scopeRepo) GetAgentIDsInLftRgtRange(minLft, maxRgt uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&models.AgentNode{}).
		Where("lft >= ? AND rgt <= ?", minLft, maxRgt).
		Pluck("agent_id", &ids).Error
	return ids, err
}

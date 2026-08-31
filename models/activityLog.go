package models

import (
	"example/go_backoffice/enums"

	"gorm.io/gorm"
)

// ActivityLog registra le azioni rilevanti eseguite nel sistema (audit trail).
// Viene popolato tipicamente dai service dopo un'operazione andata a buon fine
// (creazione, modifica, cambio stato, eliminazione, login, ecc.).
type ActivityLog struct {
	gorm.Model

	// Chi ha eseguito l'azione (può essere nil per eventi di sistema/anonimi,
	// es. tentativo di login fallito prima dell'autenticazione)
	ActorID *uint `json:"actor_id" gorm:"index"`
	Actor   *User `json:"actor,omitempty" gorm:"foreignKey:ActorID"`

	// Ruolo dell'attore al momento dell'azione (snapshot: il ruolo potrebbe
	// cambiare in futuro, ma il log deve restare storicamente accurato)
	ActorRole enums.Role `json:"role" gorm:"type:enum('ADMIN', 'OPERATOR', 'AGENT', 'AGENCY', 'USER')"`

	// Tipo di azione compiuta
	Action enums.Action `json:"action" gorm:"type:enum('CREATE', 'UPDATE', 'DELETE', 'RESTORE', 'ASSIGNMENT', 'REMOVE', 'MOVE', 'ACTIVE', 'SUSPEND', 'BLOCK')"`

	// Entità target su cui è stata eseguita l'azione (es. "User", "AgentNode",
	// "AgentOperator"). Vuoto per azioni senza un target diretto (es. LOGIN)
	TargetType string `json:"target_type" gorm:"type:varchar(50);index"`

	// ID dell'entità target, nullable per azioni senza target
	TargetID *uint `json:"target_id" gorm:"index"`

	// Descrizione testuale leggibile dell'evento (es. "Operatore #3 ha
	// sospeso l'utente #12")
	Description string `json:"description" gorm:"type:text"`

	// Dettagli aggiuntivi in formato JSON (es. old/new status, diff dei campi
	// modificati), serializzati come stringa per restare agnostici rispetto
	// al driver DB
	Metadata string `json:"metadata,omitempty" gorm:"type:text"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}

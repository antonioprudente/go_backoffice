package enums

import (
	"database/sql/driver"
	"fmt"
)

type Action string

const (
	Create     Action = "CREATE"
	Update     Action = "UPDATE"
	Delete     Action = "DELETE"
	Restore    Action = "RESTORE"
	Assignment Action = "ASSIGNMENT"
	Remove     Action = "REMOVE"
	Move       Action = "MOVE"
	Active     Action = "ACTIVE"
	Suspend    Action = "SUSPEND"
	Block      Action = "BLOCK"
)

func (s Action) String() string {
	return string(s)
}

// ParseAction converte una stringa nel tipo Action, validando che il valore
// corrisponda a un'azione ammessa.
func ParseAction(value string) (Action, error) {
	switch Action(value) {
	case Create, Update, Delete, Restore, Assignment, Remove, Active, Suspend, Block:
		return Action(value), nil
	default:
		return "", fmt.Errorf("azione non valida: %q", value)
	}
}

// Value implementa driver.Valuer per la conversione da Go a SQL
func (s Action) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implementa sql.Scanner per la conversione da SQL a Go
func (s *Action) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*s = Action(string(v))
	case string:
		*s = Action(v)
	default:
		return fmt.Errorf("impossibile convertire %T in Action", value)
	}
	return nil
}

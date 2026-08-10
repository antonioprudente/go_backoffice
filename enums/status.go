package enums

import (
	"database/sql/driver"
	"fmt"
)

// 1. Definisci il tipo custom
type Status string

// 2. Costanti riutilizzabili ovunque nel progetto
const (
	StatusActive    Status = "ACTIVE"
	StatusSuspended Status = "SUSPENDED"
	StatusBlocked   Status = "BLOCKED"
	StatusDefault   Status = "DEFAULT"
)

// 3. (Opzionale ma raccomandato per database) Implementazione di driver.Valuer e sql.Scanner
// Permettono a GORM e database/sql di gestire la conversione in automatico ed in modo sicuro.

func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *Status) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("impossibile convertire %v in Status", value)
		}
		*s = Status(str)
		return nil
	}
	*s = Status(string(bytes))
	return nil
}

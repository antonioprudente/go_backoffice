package enums

import (
	"database/sql/driver"
	"fmt"
)

type Role string

const (
	RoleAdmin    Role = "ADMIN"
	RoleOperator Role = "OPERATOR"
	RoleAgent    Role = "AGENT"
	RoleAgency   Role = "AGENCY"
	RoleUser     Role = "USER"
)

func (s Role) String() string {
	return string(s)
}

// ParseRole converte una stringa nel tipo Role, validando che il valore
// corrisponda a uno dei ruoli ammessi. Ritorna un errore se la stringa
// non corrisponde a nessun ruolo valido.
func ParseRole(value string) (Role, error) {
	switch Role(value) {
	case RoleAdmin, RoleOperator, RoleAgent, RoleAgency, RoleUser:
		return Role(value), nil
	default:
		return "", fmt.Errorf("ruolo non valido: %q", value)
	}
}

func (s Role) RoleValue() (driver.Value, error) {
	return string(s), nil
}

func (s *Role) ScanRole(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("impossibile convertire %v in Role", value)
		}
		*s = Role(str)
		return nil
	}
	*s = Role(string(bytes))
	return nil
}

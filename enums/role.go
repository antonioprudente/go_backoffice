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

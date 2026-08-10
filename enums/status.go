package enums

import (
	"database/sql/driver"
	"fmt"
)

type Status string

const (
	StatusActive    Status = "ACTIVE"
	StatusSuspended Status = "SUSPENDED"
	StatusBlocked   Status = "BLOCKED"
	StatusDefault   Status = "DEFAULT"
)

func (s Status) StatusValue() (driver.Value, error) {
	return string(s), nil
}

func (s *Status) ScanStatus(value interface{}) error {
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

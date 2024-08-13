package sqlc

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type Decimal float64

// Value implements the driver.Valuer interface.
func (d Decimal) Value() (driver.Value, error) {
	return float64(d), nil
}

// Scan implements the sql.Scanner interface.
func (d *Decimal) Scan(value interface{}) error {
	if value == nil {
		*d = 0
		return nil
	}

	switch v := value.(type) {
	case float64:
		*d = Decimal(v)
		return nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("failed to parse decimal: %v", err)
		}
		*d = Decimal(f)
		return nil
	default:
		return fmt.Errorf("unsupported type for Decimal: %T", value)
	}
}

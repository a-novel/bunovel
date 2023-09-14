package bunovel

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

// Duration extends the time.Duration type to be readable by a SQL driver.
// https://bun.uptrace.dev/guide/custom-types.html#sql-scanner
type Duration time.Duration

var _ sql.Scanner = (*Duration)(nil)
var _ driver.Valuer = (*Duration)(nil)

func (duration *Duration) Scan(src interface{}) error {
	switch tsrc := src.(type) {
	case string:
		value, err := time.ParseDuration(tsrc)
		*duration = Duration(value)
		return err
	case []byte:
		value, err := time.ParseDuration(string(tsrc))
		*duration = Duration(value)
		return err
	case nil:
		*duration = Duration(0)
		return nil
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}

func (duration Duration) Value() (driver.Value, error) {
	// In SQL, a duration is best represented as an int64 (or BIGINT) number.
	return int64(duration), nil
}

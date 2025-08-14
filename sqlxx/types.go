package sqlxx

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// swagger:type string
// swagger:model nullString
type NullString string

// NullTime implements sql.NullTime functionality.
//
// swagger:model nullTime
// required: false
type NullTime time.Time

// Scan implements the Scanner interface.
func (ns *NullTime) Scan(value interface{}) error {
	var v sql.NullTime
	if err := (&v).Scan(value); err != nil {
		return err
	}
	*ns = NullTime(v.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullTime) Value() (driver.Value, error) {
	return sql.NullTime{Valid: !time.Time(ns).IsZero(), Time: time.Time(ns)}.Value()
}

// JSONRawMessage represents a json.RawMessage that works well with JSON, SQL, and Swagger.
type JSONRawMessage json.RawMessage

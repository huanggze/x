package sqlxx

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

// JSONScan is a generic helper for storing a value as a JSON blob in SQL.
func JSONScan(dst interface{}, value interface{}) error {
	if value == nil {
		value = "null"
	}
	if err := json.Unmarshal([]byte(fmt.Sprintf("%s", value)), &dst); err != nil {
		return fmt.Errorf("unable to decode payload to: %s", err)
	}
	return nil
}

// JSONValue is a generic helper for retrieving a SQL JSON-encoded value.
func JSONValue(src interface{}) (driver.Value, error) {
	if src == nil {
		return nil, nil
	}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&src); err != nil {
		return nil, err
	}
	return b.String(), nil
}

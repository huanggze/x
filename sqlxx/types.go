package sqlxx

import (
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

// JSONRawMessage represents a json.RawMessage that works well with JSON, SQL, and Swagger.
type JSONRawMessage json.RawMessage

package sqlxx

import "encoding/json"

// swagger:type string
// swagger:model nullString
type NullString string

// JSONRawMessage represents a json.RawMessage that works well with JSON, SQL, and Swagger.
type JSONRawMessage json.RawMessage

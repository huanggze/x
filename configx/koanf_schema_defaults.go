package configx

import (
	"errors"
	"strings"

	"github.com/knadh/koanf/maps"

	"github.com/huanggze/x/jsonschemax"
	"github.com/ory/jsonschema/v3"
)

type KoanfSchemaDefaults struct {
	keys []jsonschemax.Path
}

func NewKoanfSchemaDefaults(rawSchema []byte, schema *jsonschema.Schema) (*KoanfSchemaDefaults, error) {
	keys, err := getSchemaPaths(rawSchema, schema)
	if err != nil {
		return nil, err
	}

	return &KoanfSchemaDefaults{keys: keys}, nil
}

func (k *KoanfSchemaDefaults) ReadBytes() ([]byte, error) {
	return nil, errors.New("schema defaults provider does not support this method")
}

func (k *KoanfSchemaDefaults) Read() (map[string]interface{}, error) {
	values := map[string]interface{}{}
	for _, key := range k.keys {
		// It's an array!
		if strings.Contains(key.Name, "#") {
			continue
		}

		if key.Default != nil {
			values[key.Name] = key.Default
		}
	}

	return maps.Unflatten(values, "."), nil
}

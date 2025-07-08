package otelx

import (
	"bytes"
	_ "embed"
	"io"
)

type ProvidersConfig struct {
}

type Config struct {
	ServiceName           string          `json:"service_name"`
	DeploymentEnvironment string          `json:"deployment_environment"`
	Provider              string          `json:"provider"`
	Providers             ProvidersConfig `json:"providers"`
}

//go:embed config.schema.json
var ConfigSchema string

const ConfigSchemaID = "ory://tracing-config"

// AddConfigSchema adds the tracing schema to the compiler.
// The interface is specified instead of `jsonschema.Compiler` to allow the use of any jsonschema library fork or version.
func AddConfigSchema(c interface {
	AddResource(url string, r io.Reader) error
}) error {
	return c.AddResource(ConfigSchemaID, bytes.NewBufferString(ConfigSchema))
}

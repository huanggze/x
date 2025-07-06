package otelx

type ProvidersConfig struct {
}

type Config struct {
	ServiceName           string          `json:"service_name"`
	DeploymentEnvironment string          `json:"deployment_environment"`
	Provider              string          `json:"provider"`
	Providers             ProvidersConfig `json:"providers"`
}

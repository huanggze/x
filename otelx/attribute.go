package otelx

import "go.opentelemetry.io/otel/attribute"

func StringAttrs(attrs map[string]string) []attribute.KeyValue {
	s := make([]attribute.KeyValue, 0, len(attrs))
	for k, v := range attrs {
		s = append(s, attribute.String(k, v))
	}
	return s
}

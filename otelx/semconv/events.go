package semconv

import (
	"github.com/gofrs/uuid"
	otelattr "go.opentelemetry.io/otel/attribute"

	"github.com/huanggze/x/httpx"
)

type Event string

func (e Event) String() string {
	return string(e)
}

type AttributeKey string

func (a AttributeKey) String() string {
	return string(a)
}

const (
	AttributeKeyIdentityID         AttributeKey = "IdentityID"
	AttributeKeyNID                AttributeKey = "ProjectID"
	AttributeKeyClientIP           AttributeKey = "ClientIP"
	AttributeKeyGeoLocationCity    AttributeKey = "GeoLocationCity"
	AttributeKeyGeoLocationRegion  AttributeKey = "GeoLocationRegion"
	AttributeKeyGeoLocationCountry AttributeKey = "GeoLocationCountry"
	AttributeKeyWorkspace          AttributeKey = "WorkspaceID"
	AttributeKeySubscriptionID     AttributeKey = "SubscriptionID"
	AttributeKeyProjectEnvironment AttributeKey = "ProjectEnvironment"
	AttributeKeyWorkspaceAPIKeyID  AttributeKey = "WorkspaceAPIKeyID"
	AttributeKeyProjectAPIKeyID    AttributeKey = "ProjectAPIKeyID"
)

func AttrIdentityID[V string | uuid.UUID](val V) otelattr.KeyValue {
	return otelattr.String(AttributeKeyIdentityID.String(), uuidOrString(val))
}

func AttrClientIP(val string) otelattr.KeyValue {
	return otelattr.String(AttributeKeyClientIP.String(), val)
}

func AttrGeoLocation(val httpx.GeoLocation) []otelattr.KeyValue {
	geoLocationAttributes := make([]otelattr.KeyValue, 0, 3)

	if val.City != "" {
		geoLocationAttributes = append(geoLocationAttributes, otelattr.String(AttributeKeyGeoLocationCity.String(), val.City))
	}
	if val.Country != "" {
		geoLocationAttributes = append(geoLocationAttributes, otelattr.String(AttributeKeyGeoLocationCountry.String(), val.Country))
	}
	if val.Region != "" {
		geoLocationAttributes = append(geoLocationAttributes, otelattr.String(AttributeKeyGeoLocationRegion.String(), val.Region))
	}

	return geoLocationAttributes
}

func uuidOrString[V string | uuid.UUID](val V) string {
	switch val := any(val).(type) {
	case string:
		return val
	case uuid.UUID:
		return val.String()
	}
	panic("unreachable")
}

package semconv

import (
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
	AttributeKeyClientIP           AttributeKey = "ClientIP"
	AttributeKeyGeoLocationCity    AttributeKey = "GeoLocationCity"
	AttributeKeyGeoLocationRegion  AttributeKey = "GeoLocationRegion"
	AttributeKeyGeoLocationCountry AttributeKey = "GeoLocationCountry"
)

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

package semconv

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"

	"github.com/huanggze/x/httpx"
)

type contextKey int

const contextKeyAttributes contextKey = iota

func ContextWithAttributes(ctx context.Context, attrs ...attribute.KeyValue) context.Context {
	existing, _ := ctx.Value(contextKeyAttributes).([]attribute.KeyValue)
	return context.WithValue(ctx, contextKeyAttributes, append(existing, attrs...))
}

func Middleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := ContextWithAttributes(r.Context(),
		append(
			AttrGeoLocation(*httpx.ClientGeoLocation(r)),
			AttrClientIP(httpx.ClientIP(r)),
		)...,
	)

	next(rw, r.WithContext(ctx))
}

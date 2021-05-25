package divert

import (
	"context"
	"net/http"
)

type divertHeaderKey string

const (
	// DivertHeaderName is the header used for divert values.
	DivertHeaderName = "x-okteto-dvrt"

	// divertHeaderCtxKey is the unique key value for divert header
	// value injected into context.
	divertHeaderCtxKey = divertHeaderKey(DivertHeaderName)
)

// FromContext provides the divert header values stored in context.
func FromContext(ctx context.Context) string {
	if v := ctx.Value(divertHeaderCtxKey); v != nil {
		val, _ := v.(string)
		return val
	}
	return ""
}

// FromHeaders extracts divert headers from an http request
// and provides the value. If missing then empty string
// is provided.
func FromHeaders(r *http.Request) string {
	return r.Header.Get(DivertHeaderName)
}

// InjectDivertHeader is an http middleware handler that
// injects Okteto divert headers into context from http.Request.
func InjectDivertHeader() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), divertHeaderCtxKey, r.Header.Get(DivertHeaderName))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

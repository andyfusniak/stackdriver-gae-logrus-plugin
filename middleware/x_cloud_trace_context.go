package middleware

import (
	"context"
	"net/http"
)

type xctcKeyType string

const xctcKey xctcKeyType = "xctc"

// XCTC returns the XCloudTraceContent value from the context.
func XCTC(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	xctc, ok := ctx.Value(xctcKey).(string)
	if !ok {
		return ""
	}
	return xctc
}

// XCloudTraceContext middleware extracts the X-Cloud-Trace-Context
// from the request header and injects it into the context. The value
// is read by the logrus GAEStandardFormatter to thread log entries
// by request.
func XCloudTraceContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), xctcKey, r.Header.Get("X-Cloud-Trace-Context"))
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

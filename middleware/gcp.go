package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/toga4/go-api-challange/log"
)

type contextKey string

const (
	traceContextKey = contextKey("gcp-trace-id")

	cloudTraceContextHeaderName = "X-Cloud-Trace-Context"
	traceContextHeaderName      = "X-Trace-Context"
)

func GCPTraceLogger(projectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			traceID := extractTraceID(r)

			if traceID != "" {
				trace := fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)
				r = r.WithContext(context.WithValue(r.Context(), traceContextKey, traceID))

				// https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
				l := log.R(r).WithValues("logging.googleapis.com/trace", trace)
				r = r.WithContext(log.NewContext(r.Context(), l))
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func extractTraceID(r *http.Request) string {
	// https://cloud.google.com/load-balancing/docs/https?hl=ja#target-proxies
	traceContext := r.Header.Get(cloudTraceContextHeaderName)
	i := strings.Index(traceContext, "/")
	if i >= 0 {
		return traceContext[:i]
	}

	traceContext = r.Header.Get(traceContextHeaderName)
	if i < 0 {
		return traceContext
	}

	return ""
}

type GCPTraceTransport struct {
	Transport http.RoundTripper
}

func (t *GCPTraceTransport) transport() http.RoundTripper {
	if t.Transport == nil {
		return http.DefaultTransport
	}
	return t.Transport
}

func (t *GCPTraceTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	trace, ok := r.Context().Value(traceContextKey).(string)
	if ok {
		r.Header.Set(traceContextHeaderName, trace)
	}
	return t.transport().RoundTrip(r)
}

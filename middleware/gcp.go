package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/toga4/go-api-challange/log"
)

func GCPTraceLogger(projectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			traceID := extractTraceID(r)

			if traceID != "" {
				// https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
				trace := fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)
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
	traceContext := r.Header.Get("X-Cloud-Trace-Context")
	i := strings.Index(traceContext, "/")
	if i < 0 {
		return ""
	}
	return traceContext[:i]
}

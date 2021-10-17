package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-logr/logr"
	"github.com/toga4/go-api-challange/log"
)

func WithLogger(l logr.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(log.NewContext(r.Context(), l)))
		}
		return http.HandlerFunc(fn)
	}
}

func RequestLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := log.R(r).WithValues(
			"host", r.Host,
			"request_uri", r.RequestURI,
			"remote_addr", r.RemoteAddr,
		)
		log.Info("request received")

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()
		defer func() {
			log.WithValues(
				"status", ww.Status(),
				"bytes_written", ww.BytesWritten(),
				"elapsed", time.Since(t1),
			).Info("response replied")
		}()

		next.ServeHTTP(ww, r)
	}

	return http.HandlerFunc(fn)
}

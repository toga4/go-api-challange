package log

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

type contextKey string

const (
	loggerContextKey = contextKey("logger")
)

func NewLogger() (logr.Logger, error) {
	z, err := zap.NewDevelopment()
	if err != nil {
		return logr.Discard(), err
	}

	return zapr.NewLogger(z), nil
}

func NewContext(ctx context.Context, l logr.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}

func R(r *http.Request, keysAndValues ...interface{}) logr.Logger {
	return C(r.Context())
}

func C(ctx context.Context, keysAndValues ...interface{}) logr.Logger {
	l, ok := ctx.Value(loggerContextKey).(logr.Logger)
	if !ok {
		panic("could not find logr.Logger from context")
	}
	return l.WithValues(keysAndValues...)
}

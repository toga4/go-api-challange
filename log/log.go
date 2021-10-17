package log

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const (
	loggerContextKey = contextKey("logger")
)

func NewLoggerForLocal() (logr.Logger, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return logr.Discard(), err
	}
	return zapr.NewLogger(l), nil
}

func NewLogger() (logr.Logger, error) {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "severity",
			NameKey:        "logger",
			MessageKey:     "message",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	l, err := config.Build()
	if err != nil {
		return logr.Discard(), err
	}

	return zapr.NewLogger(l), nil
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

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-logr/logr"
	"github.com/kelseyhightower/envconfig"
	"github.com/toga4/go-api-challange/interfaces/handler"
	"github.com/toga4/go-api-challange/log"
	"github.com/toga4/go-api-challange/middleware"
	"github.com/toga4/go-api-challange/usecase"
)

type Env struct {
	Env          string `envconfig:"GO_ENV" default:"local"`
	Port         string `envconfig:"PORT" default:"8000"`
	GCPProjectID string `envconfig:"GCP_PROJECT_ID"`
	HostURI      string `envconfig:"HOST_URI" default:"http://localhost:8000"`
}

func main() {
	var env Env
	if err := envconfig.Process("", &env); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse environment variables: %v\n", err.Error())
		os.Exit(1)
	}

	// Initialize logger
	logger, err := newLogger(env.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err.Error())
		os.Exit(1)
	}

	// Dependency Injection
	cu := usecase.NewChallangeUsecase(env.HostURI)
	ch := handler.NewChallangeHandler(cu)

	// Routing
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(middleware.WithLogger(logger))
	r.Use(middleware.GCPTraceLogger(env.GCPProjectID))
	r.Use(middleware.RequestLogger)
	r.Use(chimw.Recoverer)
	r.Get("/healthz", ch.HandleHealthCheck)
	r.Get("/", ch.HandleHello)
	r.Get("/delegate", ch.HandleDelegate)

	// Initialize Server
	server := &http.Server{Addr: ":" + env.Port, Handler: r}

	// Graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	// Start server
	fmt.Fprintf(os.Stderr, "Server listening on port %s.\n", env.Port)
	fmt.Fprintln(os.Stderr, server.ListenAndServe())
}

func newLogger(env string) (logr.Logger, error) {
	if env == "local" {
		return log.NewLoggerForLocal()
	} else {
		return log.NewLogger()
	}
}

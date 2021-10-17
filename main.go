package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kelseyhightower/envconfig"
	"github.com/toga4/go-api-challange/interfaces/handler"
)

type Env struct {
	Port string `envconfig:"PORT" default:"8000"`
}

func main() {
	var env Env
	if err := envconfig.Process("", &env); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse environment variables: %v\n", err.Error())
		os.Exit(1)
	}

	// Dependency Injection
	ch := handler.NewChallangeHandler()

	// Routing
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/healthz", ch.HandleHealthCheck)

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
// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	appfiles "app-name"
	"app-name/internal/api"
	"app-name/internal/apiserver"
	"app-name/internal/config"
	"app-name/internal/platform"
	"app-name/internal/store"
	"app-name/internal/worker"
)

func main() {
	if err := run(); err != nil {
		slog.Error("app stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	env, err := config.FromEnv()
	if err != nil {
		return err
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: env.LogLevel}))
	slog.SetDefault(log)
	log.Info("starting", "commit", api.GitCommit, "built", api.BuildTimestamp)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := store.OpenSQLite(ctx, env.DataDir)
	if err != nil {
		return err
	}
	defer db.Close()

	p := platform.NewClient(env.ApiEndpoint, env.ApiToken)
	if err := p.InitAssetTypes(appfiles.AssetTypes); err != nil {
		return fmt.Errorf("init asset types: %w", err)
	}

	srv, err := api.NewServer(db, appfiles.OpenAPI)
	if err != nil {
		return err
	}
	handler := apiserver.HandlerWithOptions(apiserver.NewStrictHandler(srv, nil), apiserver.StdHTTPServerOptions{BaseURL: "/v1"})
	mux := http.NewServeMux()
	mux.Handle("/v1/", handler)
	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", env.ApiServerPort), Handler: mux, ReadHeaderTimeout: 10 * time.Second}

	// probes on port 5000 like every BuildingPro Suites service: alive while the process runs, ready once the
	// asset types exist and the API is listening
	var ready atomic.Bool
	probes := http.NewServeMux()
	probes.HandleFunc("GET /livez", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	probes.HandleFunc("GET /readyz", func(w http.ResponseWriter, _ *http.Request) {
		if !ready.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	probeServer := &http.Server{Addr: ":5000", Handler: probes, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := probeServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("probe server", "err", err)
		}
	}()

	go worker.New(db, p, log).Run(ctx)

	listener, err := net.Listen("tcp", httpServer.Addr)
	if err != nil {
		return err
	}
	errc := make(chan error, 1)
	go func() {
		log.Info("API listening", "port", env.ApiServerPort)
		if err := httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errc <- err
		}
		close(errc)
	}()
	ready.Store(true)

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	log.Info("shutting down")
	ready.Store(false)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	probeServer.Shutdown(shutdownCtx)
	return httpServer.Shutdown(shutdownCtx)
}

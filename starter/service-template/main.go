package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/azghr/forge/envconfig"
	"github.com/azghr/forge/retry"
	"github.com/azghr/forge/stopwatch"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var cfg Config
	if err := envconfig.Load(&cfg); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	setupLogger(cfg)

	slog.Info("starting service", "config", cfg.String())

	worker := NewWorker(cfg)
	go worker.Run(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /api/tasks", handleTasks)
	mux.HandleFunc("POST /api/tasks", handleCreateTask)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      withMiddleware(mux, cfg),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("http server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var sw stopwatch.Stopwatch
	sw.Start()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
	slog.Info("server stopped", "shutdown_duration", sw.Elapsed())

	worker.Stop()
}

func setupLogger(cfg Config) {
	level := parseLogLevel(cfg.LogLevel)
	var h slog.Handler
	switch cfg.LogFormat {
	case "json":
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	default:
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
	slog.SetDefault(slog.New(h))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

var tasks []Task

func handleTasks(w http.ResponseWriter, r *http.Request) {
	var sw stopwatch.Stopwatch
	sw.Start()
	defer func() {
		slog.Debug("handleTasks", "duration", sw.Elapsed())
	}()

	respondJSON(w, http.StatusOK, tasks)
}

func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := decodeJSON(r, &task); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := task.Validate(); err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var sw stopwatch.Stopwatch
	sw.Start()
	err := retry.RetryContext(r.Context(), retry.RetryConfig{
		MaxTries:   3,
		InitDelay:  50 * time.Millisecond,
		Multiplier: 2.0,
		MaxDelay:   2 * time.Second,
	}, func() error {
		tasks = append(tasks, task)
		slog.Info("task created", "task", task)
		return nil
	})
	slog.Debug("handleCreateTask", "retry_duration", sw.Elapsed())

	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

package main

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/azghr/forge/retry"
	"github.com/azghr/forge/stopwatch"
)

type Worker struct {
	cfg    Config
	wg     sync.WaitGroup
	stopCh chan struct{}
}

func NewWorker(cfg Config) *Worker {
	return &Worker{
		cfg:    cfg,
		stopCh: make(chan struct{}),
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.wg.Add(1)
	defer w.wg.Done()

	slog.Info("worker started")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("worker stopped by context")
			return
		case <-w.stopCh:
			slog.Info("worker stopped")
			return
		case <-ticker.C:
			w.processTasks(ctx)
		}
	}
}

func (w *Worker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
}

func (w *Worker) processTasks(ctx context.Context) {
	var sw stopwatch.Stopwatch
	sw.Start()
	defer func() {
		slog.Debug("worker cycle completed", "duration", sw.Elapsed())
	}()

	err := retry.RetryContext(ctx, retry.RetryConfig{
		MaxTries:   3,
		InitDelay:  100 * time.Millisecond,
		Multiplier: 2.0,
		MaxDelay:   5 * time.Second,
	}, func() error {
		return w.doWork(ctx)
	})
	if err != nil {
		slog.Error("worker cycle failed after retries", "error", err)
	}
}

func (w *Worker) doWork(ctx context.Context) error {
	slog.Debug("processing tasks")
	return nil
}

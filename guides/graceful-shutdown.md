# Graceful Shutdown in Go

A production service must handle shutdown signals cleanly — draining in-flight requests, closing connections, flushing buffers, and releasing resources.

## The Pattern

```go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()

// Start server in background goroutine
srv := &http.Server{Addr: ":8080", Handler: mux}
go func() {
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
}()

// Block until signal received
<-ctx.Done()

// Graceful shutdown with timeout
shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
srv.Shutdown(shutdownCtx)
```

## Key Principles

1. **Signal handling first** — register `SIGINT` and `SIGTERM` before starting goroutines
2. **Block until done** — `<-ctx.Done()` keeps main alive until signal arrives
3. **Timeout the shutdown** — never wait forever; use `context.WithTimeout`
4. **Reverse order** — stop accepting new work, then drain existing work, then close resources

## Multiple Services

When running multiple servers or workers, use a `sync.WaitGroup`:

```go
var wg sync.WaitGroup

wg.Add(1)
go func() {
    defer wg.Done()
    // start server
}()

wg.Add(1)
go func() {
    defer wg.Done()
    // start worker
}()

<-ctx.Done()
// trigger shutdown for each component
wg.Wait()
```

## Forge Integration

The `starter/service-template/` demonstrates this pattern with an HTTP server and background worker.

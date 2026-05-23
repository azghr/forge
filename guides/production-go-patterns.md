# Production Go Patterns

A collection of patterns that make Go services production-ready. These complement the Forge packages and the starter/service-template.

## Configuration from Environment

Use `envconfig` to load structured configuration from environment variables:

```go
type Config struct {
    Port      string `env:"PORT,default=8080"`
    DatabaseURL string `env:"DATABASE_URL,required"`
    LogLevel  string `env:"LOG_LEVEL,default=info"`
}

var cfg Config
if err := envconfig.Load(&cfg); err != nil {
    log.Fatal(err)
}
```

Never use raw `os.Getenv` in application code. Centralize configuration into a single struct with validation.

## Structured Logging

Use `log/slog` with a consistent structure:

```go
slog.Info("request completed",
    "method", r.Method,
    "path", r.URL.Path,
    "status", status,
    "duration", sw.Elapsed(),
)
```

Add request-scoped context to every log line using middleware.

## Middleware Chain

Compose middleware as functions that wrap `http.Handler`:

```go
func withLogging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sw := stopwatch.Start()
        ww := &respWriter{ResponseWriter: w, status: http.StatusOK}
        next.ServeHTTP(ww, r)
        slog.Info("request", "method", r.Method, "path", r.URL.Path,
            "status", ww.status, "duration", sw.Elapsed())
    })
}
```

## JSON API Responses

Standardize on a consistent JSON response format:

```go
type APIResponse struct {
    Data  any    `json:"data,omitempty"`
    Error string `json:"error,omitempty"`
}

func respondJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(APIResponse{Data: data})
}

func respondError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(APIResponse{Error: msg})
}
```

## Server Timeouts

Always set timeouts on `http.Server`:

```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
}
```

This prevents slow clients from consuming server resources indefinitely.

## Resource Lifecycle

Use `defer` for resource cleanup, never rely on finalizers:

```go
f, err := os.Create(path)
if err != nil {
    return err
}
defer f.Close()
```

For complex resources, use a `Close` method with `sync.Once`:

```go
type Server struct {
    closeOnce sync.Once
    closer    io.Closer
}

func (s *Server) Close() error {
    var err error
    s.closeOnce.Do(func() {
        err = s.closer.Close()
    })
    return err
}
```

## Background Workers

Separate worker logic from HTTP handlers. Use the pattern from `starter/service-template/`:

```go
type Worker struct {
    stopCh chan struct{}
    wg     sync.WaitGroup
}

func (w *Worker) Start(ctx context.Context) {
    w.wg.Add(1)
    go func() {
        defer w.wg.Done()
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                w.processBatch(ctx)
            }
        }
    }()
}
```

## Key Rules

1. Centralize configuration — one struct, one `envconfig.Load` call
2. Log everything at the boundary — request start/end, errors, slow operations
3. Set timeouts — on HTTP servers, HTTP clients, and database connections
4. Clean up resources — every `Open` needs a corresponding `Close`
5. Separate concerns — HTTP handling, business logic, and background work should live in different files
6. Test the integration — the starter/service-template shows how to test config, handlers, and workers together

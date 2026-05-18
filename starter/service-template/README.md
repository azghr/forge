# Service Template

Production-ready Go HTTP service built with Forge packages.

## Quick Start

```sh
go run .
```

## What's Inside

| Concern | Implementation |
|---------|---------------|
| Config | `forge/envconfig` — environment-based config with defaults |
| HTTP Server | `net/http` with graceful shutdown via `signal.NotifyContext` |
| Middleware | Request logging, panic recovery, observability |
| Retry | `forge/retry` — exponential backoff with full-jitter |
| Workers | Background task processing with graceful stop |
| Observability | `forge/stopwatch` — request timing and metrics |
| Logging | `log/slog` — structured logging (text or JSON) |
| Validation | Struct validation via `forge/validator` |

## Config

| Env Var | Default | Description |
|---------|---------|-------------|
| `PORT` | `8080` | HTTP listen port |
| `LOG_FORMAT` | `text` | Log format (`text` or `json`) |
| `LOG_LEVEL` | `0` | Log level (0=info, -4=debug) |

## API

- `GET /health` — Health check
- `GET /api/tasks` — List tasks
- `POST /api/tasks` — Create task

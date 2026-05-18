package main

import (
	"log/slog"
	"os"
)

type Config struct {
	Port      string `env:"PORT,default=8080"`
	LogFormat string `env:"LOG_FORMAT,default=text"`
	LogLevel  string `env:"LOG_LEVEL,default=info"`
	DBPath    string `env:"DB_PATH,default=data/tasks.json"`
}

func (c Config) String() string {
	return "port=" + c.Port + " log_format=" + c.LogFormat + " log_level=" + c.LogLevel
}

func parseLogLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

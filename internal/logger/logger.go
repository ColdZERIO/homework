package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	Level     slog.Level
	AddSource bool
}

func New(cfg Config) *slog.Logger {
	options := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	handler := slog.NewJSONHandler(os.Stdout, options)

	return slog.New(handler)
}
package logger

import (
	"context"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
	"path/filepath"
)

type MultiHandler struct {
	handlers []slog.Handler
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}

type Config struct {
	LogDir     string
	LogFile    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

func Init(cfg Config) func() {
	fullPath := filepath.Join(cfg.LogDir, cfg.LogFile)

	lumberjackLogger := &lumberjack.Logger{
		Filename:   fullPath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}

	fileHandler := slog.NewJSONHandler(lumberjackLogger, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	multiHandler := &MultiHandler{
		handlers: []slog.Handler{consoleHandler, fileHandler},
	}

	logger := slog.New(multiHandler)
	slog.SetDefault(logger)

	return func() {
		_ = lumberjackLogger.Close()
	}
}

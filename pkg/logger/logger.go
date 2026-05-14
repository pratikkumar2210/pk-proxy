package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	*slog.Logger
}

func getLogLevelFromEnv() slog.Level {
	levelStr := os.Getenv("LOG_LEVEL")

	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func NewLogger(hostname string) *Logger {
	logFile, err := os.OpenFile("/Users/pratikkumar/Desktop/dumps/pk-proxy/logs/api.example.com", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// NOTE: Do NOT defer logFile.Close() here — the file must stay open for the lifetime of the logger.

	handler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		AddSource: false,
		Level:     getLogLevelFromEnv(),
	})
	return &Logger{
		slog.New(handler),
	}
}

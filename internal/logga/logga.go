package logga

import (
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/test/testbuffer"
)

// ErrInvalidLogWriter and ErrInvalidLogLevelValue are errors returned when the logger config is invalid.
var (
	ErrInvalidLogWriter     = errors.New("a valid writer was not defined in configuration")
	ErrInvalidLogLevelValue = errors.New("a valid log level was not defined in configuration")
)

// Logga wraps slog.Logger to provide structured, context-aware logging.
type Logga struct {
	Log *slog.Logger
}

// KVPs is a map of key-value pairs passed as extra context to log calls.
type KVPs map[string]any

// New will return the log handler with the options defined in config.
func New(cfg *config.Logging) (*Logga, error) {
	level, err := getLogLevelFromConfig(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	logLevel := &slog.LevelVar{}
	logLevel.Set(level)

	writer, err := getWriterFromConfig(cfg.Writer)
	if err != nil {
		return nil, err
	}

	var logger *slog.Logger

	switch cfg.OutputFormat {
	case "json":
		logger = slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: logLevel}))
	case "text":
		logger = slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{Level: logLevel}))
	default:
		logger = slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{Level: logLevel}))
	}

	return &Logga{Log: logger}, nil
}

func getWriterFromConfig(configuredWriter string) (io.Writer, error) {
	switch configuredWriter {
	case "stdout":
		return os.Stdout, nil
	case "": // Default to stdout if none is defined.
		return os.Stdout, nil
	case "testwriter": // Used for testing.
		return testbuffer.TestBuffer{}, nil
	}

	return nil, ErrInvalidLogWriter
}

func getLogLevelFromConfig(configuredLevel string) (slog.Level, error) {
	switch configuredLevel {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	}

	return slog.LevelInfo, ErrInvalidLogLevelValue
}

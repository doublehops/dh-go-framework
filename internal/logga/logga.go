package logga

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/test/testbuffer"
)

const (
	logFormatJSON   = "json"
	logWriterStdout = "stdout"
	logWriterTest   = "testwriter"
	logLevelDebug   = "DEBUG"
	logLevelInfo    = "INFO"
	logLevelWarn    = "WARN"
	logLevelError   = "ERROR"
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

	opts := &slog.HandlerOptions{Level: logLevel, ReplaceAttr: lowercaseLevel}

	switch cfg.OutputFormat {
	case logFormatJSON:
		logger = slog.New(slog.NewJSONHandler(writer, opts))
	case "text":
		logger = slog.New(slog.NewTextHandler(writer, opts))
	default:
		logger = slog.New(slog.NewTextHandler(writer, opts))
	}

	return &Logga{Log: logger}, nil
}

func getWriterFromConfig(configuredWriter string) (io.Writer, error) {
	switch configuredWriter {
	case logWriterStdout:
		return os.Stdout, nil
	case "": // Default to stdout if none is defined.
		return os.Stdout, nil
	case logWriterTest: // Used for testing.
		return testbuffer.TestBuffer{}, nil
	}

	return nil, ErrInvalidLogWriter
}

// getLogLevelFromConfig maps a configured level name onto a slog.Level.
//
// Matching is case-insensitive on purpose. The value arrives from an environment variable, an
// unrecognised one is a fatal startup error, and "info" versus "INFO" is not a distinction anyone
// should have to remember at deploy time.
func getLogLevelFromConfig(configuredLevel string) (slog.Level, error) {
	switch strings.ToUpper(strings.TrimSpace(configuredLevel)) {
	case logLevelDebug:
		return slog.LevelDebug, nil
	case logLevelInfo:
		return slog.LevelInfo, nil
	case logLevelWarn:
		return slog.LevelWarn, nil
	case logLevelError:
		return slog.LevelError, nil
	}

	return slog.LevelInfo, ErrInvalidLogLevelValue
}

// lowercaseLevel renders the level as "info" rather than slog default "INFO".
//
// Lowercase is the prevailing convention for structured logs - zerolog, zap and logrus all emit
// it, and Elastic Common Schema specifies log.level in lowercase - which matters because these
// logs are shipped to CloudWatch where consistency with other tooling makes filtering easier.
// slog is the outlier here, so this normalises it rather than the other way round.
func lowercaseLevel(_ []string, a slog.Attr) slog.Attr {
	if a.Key != slog.LevelKey {
		return a
	}

	if level, ok := a.Value.Any().(slog.Level); ok {
		a.Value = slog.StringValue(strings.ToLower(level.String()))
	}

	return a
}

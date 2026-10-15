package logga

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/test/testbuffer"
)

const (
	logFormatJSON   = "json"
	logWriterStdout = "stdout"
	logWriterFile   = "file"
	logWriterTest   = "testwriter"
	logLevelDebug   = "DEBUG"
	logLevelInfo    = "INFO"
	logLevelWarn    = "WARN"
	logLevelError   = "ERROR"
)

// ErrInvalidLogWriter, ErrMissingLogFilePath and ErrInvalidLogLevelValue are errors returned when the logger config is invalid.
var (
	ErrInvalidLogWriter     = errors.New("a valid writer was not defined in configuration")
	ErrMissingLogFilePath   = errors.New("log writer is file but no filePath was defined in configuration")
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

	writer, err := getWriterFromConfig(cfg.Writer, cfg.FilePath)
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

func getWriterFromConfig(configuredWriter, filePath string) (io.Writer, error) {
	switch configuredWriter {
	case logWriterStdout:
		return os.Stdout, nil
	case logWriterFile:
		f, err := openLogFile(filePath)
		if err != nil {
			return nil, err
		}

		// Still write to stdout so the console (and anything collecting container output) sees
		// the same log lines as the file.
		return io.MultiWriter(os.Stdout, f), nil
	case "": // Default to stdout if none is defined.
		return os.Stdout, nil
	case logWriterTest: // Used for testing.
		return testbuffer.TestBuffer{}, nil
	}

	return nil, ErrInvalidLogWriter
}

// openLogFile opens filePath for appending, creating it and its directory if needed. The handle is
// left open for the life of the process.
func openLogFile(filePath string) (*os.File, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, ErrMissingLogFilePath
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0o750); err != nil {
		return nil, fmt.Errorf("unable to create log directory. %w", err)
	}

	f, err := os.OpenFile(filepath.Clean(filePath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("unable to open log file `%s`. %w", filePath, err)
	}

	return f, nil
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

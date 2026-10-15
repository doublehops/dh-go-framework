package logga

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/doublehops/dh-go-framework/internal/app"
	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/test/testbuffer"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		config        *config.Logging
		expectedError error
	}{
		{
			name: "success",
			config: &config.Logging{
				OutputFormat: "text",
				LogLevel:     "DEBUG",
			},
			expectedError: nil,
		},
		{
			name: "success",
			config: &config.Logging{
				OutputFormat: "json",
				LogLevel:     "DEBUG",
			},
			expectedError: nil,
		},
		{
			name: "success",
			config: &config.Logging{
				OutputFormat: "",
				LogLevel:     "DEBUG",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.config)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("error not as expected. Wanted: %s; got: %s", tt.expectedError, err)
			}
		})
	}
}

func TestGetWriterFromConfig(t *testing.T) {
	tests := []struct {
		name          string
		writer        string
		filePath      string
		expectedError error
	}{
		{name: "successStdout", writer: "stdout"},
		{name: "successEmptyDefaultsToStdout", writer: ""},
		{name: "successFile", writer: "file", filePath: "app.log"},
		{name: "successFileCreatesDirectory", writer: "file", filePath: "nested/dir/app.log"},
		{name: "failFileWithoutPath", writer: "file", filePath: "", expectedError: ErrMissingLogFilePath},
		{name: "failUnknownWriter", writer: "syslog", expectedError: ErrInvalidLogWriter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.filePath
			if filePath != "" {
				filePath = filepath.Join(t.TempDir(), filePath)
			}

			_, err := getWriterFromConfig(tt.writer, filePath)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("error not as expected. Wanted: %s; got: %s", tt.expectedError, err)
			}

			if err == nil && filePath != "" {
				if _, err := os.Stat(filePath); err != nil {
					t.Errorf("expected log file to be created. %s", err)
				}
			}
		})
	}
}

// The file writer must still write to stdout as well as to the file.
func TestFileWriterAlsoWritesToStdout(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "app.log")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("unable to create pipe. %s", err)
	}

	origStdout := os.Stdout
	os.Stdout = w

	l, err := New(&config.Logging{Writer: "file", FilePath: filePath, OutputFormat: "json", LogLevel: "DEBUG"})

	os.Stdout = origStdout

	if err != nil {
		t.Fatalf("Got unexpected error. %s", err)
	}

	l.Info(context.Background(), basicMessage, nil)
	_ = w.Close()

	stdout, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("unable to read stdout. %s", err)
	}

	fileData, err := os.ReadFile(filePath) //nolint:gosec
	if err != nil {
		t.Fatalf("unable to read log file. %s", err)
	}

	for name, out := range map[string]string{"stdout": string(stdout), "file": string(fileData)} {
		if !strings.Contains(out, basicMessage) {
			t.Errorf("expected %s to contain the log message; got: %q", name, out)
		}
	}
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name            string
		configuredLevel string
		expectedLevel   slog.Level
		expectedError   error
	}{
		{
			name:            "successDEBUG",
			configuredLevel: "DEBUG",
			expectedLevel:   slog.LevelDebug,
			expectedError:   nil,
		},
		{
			name:            "successINFO",
			configuredLevel: "INFO",
			expectedLevel:   slog.LevelInfo,
			expectedError:   nil,
		},
		{
			// Parsing is case-insensitive: the value comes from an env var and an unrecognised
			// one is a fatal startup error, so casing must not decide whether the app boots.
			name:            "successLowercase",
			configuredLevel: "info",
			expectedLevel:   slog.LevelInfo,
			expectedError:   nil,
		},
		{
			name:            "successMixedCaseAndWhitespace",
			configuredLevel: "  Debug ",
			expectedLevel:   slog.LevelDebug,
			expectedError:   nil,
		},
		{
			name:            "failUnrecognisedValueStillRejected",
			configuredLevel: "verbose",
			expectedLevel:   slog.LevelInfo,
			expectedError:   ErrInvalidLogLevelValue,
		},
		{
			name:            "successWARN",
			configuredLevel: "WARN",
			expectedLevel:   slog.LevelWarn,
			expectedError:   nil,
		},
		{
			name:            "successERROR",
			configuredLevel: "ERROR",
			expectedLevel:   slog.LevelError,
			expectedError:   nil,
		},
		{
			name:            "failNotConfigured",
			configuredLevel: "",
			expectedLevel:   slog.LevelInfo,
			expectedError:   ErrInvalidLogLevelValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := getLogLevelFromConfig(tt.configuredLevel)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("error not as expected. Wanted: %s; got: %s", tt.expectedError, err)
			}
			if tt.expectedLevel != level {
				t.Errorf("level not as expected. Expected: %s; got: %s", tt.expectedLevel, level)
			}
		})
	}
}

type logOutput struct {
	Level   string `json:"level"`
	Msg     string `json:"msg"`
	TraceID string `json:"traceID"`
	UserID  int    `json:"userID"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
}

var basicMessage = "New log message"

//nolint:funlen
func TestSendLogMessage(t *testing.T) {
	_ = os.Remove(testbuffer.Filename)

	tests := []struct {
		name           string
		config         *config.Logging
		ctxArgs        map[interface{}]interface{}
		customArgs     KVPs
		expectedOutput logOutput
	}{
		{
			name: "success",
			config: &config.Logging{
				Writer:       "testwriter",
				OutputFormat: "json",
				LogLevel:     "DEBUG",
			},
			ctxArgs:    map[interface{}]interface{}{},
			customArgs: KVPs{},
			expectedOutput: logOutput{
				Level: "info",
				Msg:   basicMessage,
			},
		},
		{
			name: "successWithContextVars",
			config: &config.Logging{
				Writer:       "testwriter",
				OutputFormat: "json",
				LogLevel:     "DEBUG",
			},
			ctxArgs: map[interface{}]interface{}{
				app.TraceIDKey: "trace-99876",
				app.UserIDKey:  123,
			},
			customArgs: KVPs{},
			expectedOutput: logOutput{
				Level:   "info",
				Msg:     basicMessage,
				TraceID: "trace-99876",
				UserID:  123,
			},
		},
		{
			name: "successWithCustomVars",
			config: &config.Logging{
				Writer:       "testwriter",
				OutputFormat: "json",
				LogLevel:     "DEBUG",
			},
			customArgs: KVPs{
				"Name": "JohnS",
				"Age":  33,
			},
			expectedOutput: logOutput{
				Level: "info",
				Msg:   basicMessage,
				Name:  "JohnS",
				Age:   33,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test if file exists and wait if so. A concurrent test may have created it.
			for fileExists() {
				time.Sleep(50 * time.Millisecond)
			}
			defer func() { _ = os.Remove(testbuffer.Filename) }()

			ctx := context.Background()

			for key, value := range tt.ctxArgs {
				switch key {
				case app.TraceIDKey:
					ctx = setContextItem(ctx, app.TraceIDKey, value)
				case app.UserIDKey:
					ctx = setContextItem(ctx, app.UserIDKey, value)
				}
			}

			l, err := New(tt.config)
			if err != nil {
				t.Errorf("Got unexpected error. %s", err)
			}

			l.Info(ctx, basicMessage, tt.customArgs)

			data, err := os.ReadFile(testbuffer.Filename)
			if err != nil {
				t.Errorf("Got unexpected error reading file. %s", err)
			}

			var output logOutput
			err = json.Unmarshal(data, &output)
			if err != nil {
				t.Errorf("Got unexpected error unmarshaling JSON. %s", err)
			}

			if !reflect.DeepEqual(tt.expectedOutput, output) {
				t.Errorf("log message not as expected. Expected: %v; got: %v", tt.expectedOutput, output)
			}
		})
	}
}

func setContextItem(ctx context.Context, key app.ContextVar, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

func fileExists() bool {
	_, err := os.Stat(testbuffer.Filename)
	if os.IsNotExist(err) {
		return false
	}

	return !os.IsNotExist(err)
}

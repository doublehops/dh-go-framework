package logga

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"reflect"
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

func TestSendLogMessage(t *testing.T) {
	os.Remove(testbuffer.Filename)

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
				Level: "INFO",
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
				Level:   "INFO",
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
				Level: "INFO",
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
			defer os.Remove(testbuffer.Filename)

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

func getContextItem(ctx context.Context, key app.ContextVar) any {
	if value := ctx.Value(key); value != nil {
		return value
	}

	return ""
}

func fileExists() bool {
	_, err := os.Stat(testbuffer.Filename)
	if os.IsNotExist(err) {
		return false
	}

	return !os.IsNotExist(err)
}

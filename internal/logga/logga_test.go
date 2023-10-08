package logga

import (
	"encoding/json"
	"errors"
	"github.com/doublehops/dhapi-example/internal/config"
	"github.com/doublehops/dhapi-example/test/testbuffer"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"
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
			expectedError:   InvalidLogLevelValue,
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
		ctxArgs        map[string]interface{}
		customArgs     []interface{}
		expectedOutput logOutput
	}{
		{
			name: "success",
			config: &config.Logging{
				Writer:       "testwriter",
				OutputFormat: "json",
				LogLevel:     "DEBUG",
			},
			ctxArgs: map[string]interface{}{},
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
			ctxArgs: map[string]interface{}{
				"traceID": "ABCD-1234",
				"userID":  2134,
			},
			expectedOutput: logOutput{
				Level:   "INFO",
				Msg:     basicMessage,
				TraceID: "ABCD-1234",
				UserID:  2134,
			},
		},
		{
			name: "successWithCustomVars",
			config: &config.Logging{
				Writer:       "testwriter",
				OutputFormat: "json",
				LogLevel:     "DEBUG",
			},
			customArgs: []interface{}{
				"Name", "JohnS",
				"Age", 33,
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

			ctx := &gin.Context{}

			for key, value := range tt.ctxArgs {
				switch key {
				case "traceID":
					ctx.Set("traceID", value)
				case "userID":
					ctx.Set("userID", value)
				}
			}

			l, err := New(tt.config)
			if err != nil {
				t.Errorf("Got unexpected error. %s", err)
			}

			l.Info(ctx, basicMessage, tt.customArgs...)

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
				t.Errorf("level not as expected. Expected: %v; got: %v", tt.expectedOutput, output)
			}
		})
	}
}

func fileExists() bool {
	_, err := os.Stat(testbuffer.Filename)
	if os.IsNotExist(err) {
		return false
	}
	return !os.IsNotExist(err)
}

package logga

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/doublehops/dh-go-framework/internal/app"
)

func TestGetArguments(t *testing.T) {
	tests := []struct {
		name         string
		context      context.Context
		args         KVPs
		ctxVars      map[app.ContextVar]any
		expectedArgs []interface{}
	}{
		{
			name:    "successWithAll",
			context: context.Background(),
			args:    KVPs{"hello": "world"},
			ctxVars: map[app.ContextVar]any{app.TraceIDKey: "ABCD-1234", app.UserIDKey: 123},
			expectedArgs: []interface{}{
				slog.String("hello", "world"),
				slog.String(app.TraceIDKey.String(), "ABCD-1234"),
				slog.Int(app.UserIDKey.String(), 123),
			},
		},
		{
			name:    "SuccessOnlyArgs",
			context: context.Background(),
			args:    KVPs{"hello": "world"},
			ctxVars: map[app.ContextVar]any{},
			expectedArgs: []interface{}{
				slog.String("hello", "world"),
			},
		},
		{
			name:    "SuccessOnlyCtxArgs",
			context: context.Background(),
			args:    KVPs{},
			ctxVars: map[app.ContextVar]any{app.TraceIDKey: "ABCD-1234", app.UserIDKey: 123},
			expectedArgs: []interface{}{
				slog.String(app.TraceIDKey.String(), "ABCD-1234"),
				slog.Int(app.UserIDKey.String(), 123),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.context
			for key, value := range tt.ctxVars {
				ctx = context.WithValue(ctx, key, value)
			}
			args := addArgs(ctx, tt.args)
			assert.ElementsMatch(t, tt.expectedArgs, args, "elements not as expected")
		})
	}
}

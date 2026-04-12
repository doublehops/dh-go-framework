// Package logga provides a structured, context-aware logging wrapper around slog.
package logga

import (
	"context"
	"log/slog"

	"github.com/doublehops/dh-go-framework/internal/app"
	"github.com/doublehops/dh-go-framework/internal/tools"
)

// Debug - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Debug(ctx context.Context, msg string, kvp KVPs) {
	if kvp == nil {
		kvp = KVPs{}
	}
	kvp["func"] = tools.CurrentFunction()
	l.Log.DebugContext(ctx, msg, addArgs(ctx, kvp)...)
}

// Info - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Info(ctx context.Context, msg string, kvp KVPs) {
	if kvp == nil {
		kvp = KVPs{}
	}
	kvp["func"] = tools.CurrentFunction()
	l.Log.InfoContext(ctx, msg, addArgs(ctx, kvp)...)
}

// Warn - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Warn(ctx context.Context, msg string, kvp KVPs) {
	if kvp == nil {
		kvp = KVPs{}
	}
	kvp["func"] = tools.CurrentFunction()
	l.Log.WarnContext(ctx, msg, addArgs(ctx, kvp)...)
}

// Error - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Error(ctx context.Context, msg string, kvp KVPs) {
	if kvp == nil {
		kvp = KVPs{}
	}
	kvp["func"] = tools.CurrentFunction()
	l.Log.ErrorContext(ctx, msg, addArgs(ctx, kvp)...)
}

// addArgs will add arguments as slog.Int, slog.String, slog.Any, etc...
func addArgs(ctx context.Context, kvps KVPs) []any {
	var atts []any

	ctxArgs := getContextKVPs(ctx)
	for key, value := range ctxArgs {
		kvps[key] = value
	}

	for key, value := range kvps {
		atts = append(atts, slog.Any(key, value))
	}

	return atts
}

// getContextKVPs will check for each known context variable and add to response.
func getContextKVPs(ctx context.Context) KVPs {
	pairs := KVPs{}

	if ctx == nil {
		return pairs
	}

	if traceID := ctx.Value(app.TraceIDKey); traceID != nil {
		pairs[app.TraceIDKey.String()] = traceID
	}

	if traceID := ctx.Value(app.UserIDKey); traceID != nil {
		pairs[app.UserIDKey.String()] = traceID
	}

	return pairs
}

package logga

import (
	"context"
	"log/slog"

	"github.com/doublehops/dh-go-framework/internal/app"
	"github.com/doublehops/dh-go-framework/internal/tools"
)

// Debug - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Debug(ctx context.Context, msg string, KVP KVPs) {
	if KVP == nil {
		KVP = KVPs{}
	}
	KVP["func"] = tools.CurrentFunction()
	l.Log.DebugContext(ctx, msg, addArgs(ctx, KVP)...)
}

// Info - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Info(ctx context.Context, msg string, KVP KVPs) {
	if KVP == nil {
		KVP = KVPs{}
	}
	KVP["func"] = tools.CurrentFunction()
	l.Log.InfoContext(ctx, msg, addArgs(ctx, KVP)...)
}

// Warn - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Warn(ctx context.Context, msg string, KVP KVPs) {
	if KVP == nil {
		KVP = KVPs{}
	}
	KVP["func"] = tools.CurrentFunction()
	l.Log.WarnContext(ctx, msg, addArgs(ctx, KVP)...)
}

// Error - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Error(ctx context.Context, msg string, KVP KVPs) {
	if KVP == nil {
		KVP = KVPs{}
	}
	KVP["func"] = tools.CurrentFunction()
	l.Log.ErrorContext(ctx, msg, addArgs(ctx, KVP)...)
}

// addArgs will add arguments as slog.Int, slog.String, slog.Any, etc...
func addArgs(ctx context.Context, KVPs KVPs) []any {
	var atts []any

	ctxArgs := getContextKVPs(ctx)
	for key, value := range ctxArgs {
		KVPs[key] = value
	}

	for key, value := range KVPs {
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

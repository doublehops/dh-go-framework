package logga

import (
	"context"

	"github.com/doublehops/dh-go-framework/internal/app"
	"github.com/doublehops/dh-go-framework/internal/tools"
)

// Debug - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Debug(ctx context.Context, msg string, KVP KVPs) {
	if KVP == nil {
		KVP = KVPs{}
	}
	f := tools.CurrentFunction()
	KVP["func"] = f
	args := getArguments(ctx, KVP)
	l.Log.DebugContext(ctx, msg, args...)
}

// Info - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Info(ctx context.Context, msg string, KVP KVPs) {
	if KVP == nil {
		KVP = KVPs{}
	}
	f := tools.CurrentFunction()
	KVP["func"] = f
	args := getArguments(ctx, KVP)
	l.Log.InfoContext(ctx, msg, args...)
}

// Warn - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Warn(ctx context.Context, msg string, KVP KVPs) {
	if KVP == nil {
		KVP = KVPs{}
	}
	f := tools.CurrentFunction()
	KVP["func"] = f
	args := getArguments(ctx, KVP)
	l.Log.WarnContext(ctx, msg, args...)
}

// Error - args should be key/value pairs separated by a space. Example: "file", "migrate.go"
func (l *Logga) Error(ctx context.Context, msg string, KVP KVPs) {
	if KVP == nil {
		KVP = KVPs{}
	}
	f := tools.CurrentFunction()
	KVP["func"] = f
	args := getArguments(ctx, KVP)
	l.Log.ErrorContext(ctx, msg, args...)
}

// getArguments will get retrieve additional values from context and KVPs.
func getArguments(ctx context.Context, KVPs KVPs) []any {
	ctxArgs := getContextAtts(ctx)

	var args []any
	for key, value := range KVPs {
		args = append(args, key, value)
	}

	if len(ctxArgs) > 0 {
		args = append(args, ctxArgs...)
	}

	return args
}

// getContextAtts will retrieve known variables from context to add the log messages.
func getContextAtts(ctx context.Context) []any {
	var args []interface{}

	if ctx == nil {
		return args
	}

	if traceID := ctx.Value("traceID"); traceID != nil {
		args = append(args, "traceID", traceID)
	}

	if userID := ctx.Value(app.UserIDKey); userID != nil {
		args = append(args, app.UserIDKey.ToString(), userID)
	}

	return args
}
